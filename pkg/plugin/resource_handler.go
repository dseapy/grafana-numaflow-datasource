package plugin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dseapy/grafana-numaflow-datasource/pkg/models"
	"github.com/dseapy/grafana-numaflow-datasource/pkg/query/scenario"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type queryTypesResponse struct {
	QueryTypes []string `json:"queryTypes"`
}

type metricNamesResponse struct {
	MetricNames []string `json:"metricNames"`
}

func (d *Datasource) CallResource(_ context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req.Path == "/query-types" && req.Method == http.MethodGet {
		queryTypes := &queryTypesResponse{
			QueryTypes: []string{
				scenario.Table,
				scenario.NodeGraph,
			},
		}
		j, err := json.Marshal(queryTypes)
		if err != nil {
			backend.Logger.Error("could not marshal queryTypes JSON", "err", err)
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
			})
		}
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusOK,
			Body:   j,
		})
	} else if req.Path == "/metric-names" && req.Method == http.MethodPost {
		var qm models.QueryModel
		if d.settings == nil {
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusBadRequest,
				Body:   []byte("datasource settings nil when trying to get metric names"),
			})
		}
		if err := qm.Unmarshall(req.Body, *d.settings); err != nil {
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusBadRequest,
				Body:   []byte("could not unmarshal query JSON"),
			})
		}

		metricNames := &metricNamesResponse{}
		switch qm.RunnableQuery.ResourceType {
		case models.PipelineResourceType:
			namespacesWithPipelines, err := d.ListNamespacesWithPipelines(d.settings.Namespace)
			if err != nil {
				backend.Logger.Error("error listing namespaces with pipelines", "err", err)
				return sender.Send(&backend.CallResourceResponse{
					Status: http.StatusInternalServerError,
					Body:   []byte("error listing namespaces with pipelines"),
				})
			}
			metricNames.MetricNames = namespacesWithPipelines
		case models.VertexResourceType:
			namespacesWithVertices, err := d.ListNamespacesWithVertices(d.settings.Namespace)
			if err != nil {
				backend.Logger.Error("error listing namespaces with vertices", "err", err)
				return sender.Send(&backend.CallResourceResponse{
					Status: http.StatusInternalServerError,
					Body:   []byte("error listing namespaces with vertices"),
				})
			}
			metricNames.MetricNames = namespacesWithVertices
		case models.IsbsvcResourceType:
			namespacesWithInterStepBufferServices, err := d.ListNamespacesWithInterStepBufferServices(d.settings.Namespace)
			if err != nil {
				backend.Logger.Error("error listing namespaces with isbsvc", "err", err)
				return sender.Send(&backend.CallResourceResponse{
					Status: http.StatusInternalServerError,
					Body:   []byte("error listing namespaces with isbsvc"),
				})
			}
			metricNames.MetricNames = namespacesWithInterStepBufferServices
		default:
			backend.Logger.Error("unknown resource type, this shouldn't ever happen", "resourceType", qm.RunnableQuery.ResourceType)
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
				Body:   []byte("error listing namespaces, resource type unknown"),
			})
		}

		j, err := json.Marshal(metricNames)
		if err != nil {
			backend.Logger.Error("could not marshal metricNames JSON", "err", err)
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
				Body:   []byte(err.Error()),
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
