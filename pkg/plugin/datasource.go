package plugin

import (
	"context"
	"errors"
	"github.com/dseapy/numaflow-datasource/pkg/client"
	"github.com/dseapy/numaflow-datasource/pkg/query"
	"github.com/dseapy/numaflow-datasource/pkg/resource"
	"github.com/dseapy/numaflow-datasource/pkg/scenario"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces- only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ backend.CallResourceHandler   = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(dis backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings, err := loadSettings(dis)
	if err != nil {
		return nil, err
	}
	ns := settings.Namespace
	if !settings.Namespaced {
		ns = v1.NamespaceAll
	}
	c, err := client.NewClient(ns)
	if err != nil {
		return nil, err
	}
	return &Datasource{
		settings: settings,
		client:   c,
	}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	settings *Settings
	client   *client.Client
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Debug("QueryData called", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := runQuery(ctx, *d.settings, d.client, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *Datasource) CallResource(_ context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req.Path == resource.QueryTypesAPIPath && req.Method == resource.QueryTypesAPIMethod {
		j, err := resource.QueryTypesJson()
		if err != nil {
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
			})
		}
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusOK,
			Body:   j,
		})
	} else if req.Path == resource.MetricNamesAPIPath && req.Method == resource.MetricNamesAPIMethod {
		var q query.Query
		if d.settings == nil {
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusBadRequest,
				Body:   []byte("datasource settings nil when trying to get metric names"),
			})
		}
		if err := q.Unmarshall(req.Body); err != nil {
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusBadRequest,
				Body:   []byte(err.Error()),
			})
		}
		j, err := resource.MetricNamesJson(&q, d.client)
		if err != nil {
			return sender.Send(&backend.CallResourceResponse{
				// might be bad request still... this is good enough for now
				Status: http.StatusInternalServerError,
				// consider not passing back, but logging instead
				Body: []byte(err.Error()),
			})
		}
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusOK,
			Body:   j,
		})
	} else {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusNotFound,
		})
	}
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Debug("CheckHealth called", "request", req)

	if _, err := d.client.ListNamespacesWithInterStepBufferServices(); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}

func runQuery(_ context.Context, settings Settings, client *client.Client, dq backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}
	var q query.Query
	backend.Logger.Debug("query json %v", string(dq.JSON))
	if response.Error = q.Unmarshall(dq.JSON); response.Error != nil {
		return response
	}

	// validate
	if settings.Namespaced {
		q.RunnableQuery.Namespace = &settings.Namespace
	}
	if q.RunnableQuery.Namespace == nil {
		q.RunnableQuery.Namespace = pointer.String(v1.NamespaceAll)
	}
	if *q.RunnableQuery.Namespace == v1.NamespaceAll && q.RunnableQuery.ResourceName != "*" {
		response.Error = errors.New(`"namespace" must be provided when requesting a single pipeline, vertex, or isbsvc by name`)
		return response
	}

	// create frames
	frames, err := scenario.NewDataFrames(client, dq, q.RunnableQuery)
	if err != nil {
		backend.Logger.Error("error retrieving frames", "err", err)
		response.Error = errors.New(`error retrieving frames`)
		return response
	}
	if len(frames) == 0 {
		return response
	}
	for _, frame := range frames {
		frame.RefID = dq.RefID
	}
	response.Frames = append(response.Frames, frames...)
	return response
}
