package plugin

import (
	"context"
	"encoding/json"
	dfv1 "github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	v1 "k8s.io/api/core/v1"
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
		if err := qm.Unmarshall(req.Body); err != nil {
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusBadRequest,
				Body:   []byte(err.Error()),
			})
		}

		// validate
		if qm.RunnableQuery.Namespace == nil || *qm.RunnableQuery.Namespace == v1.NamespaceAll {
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusBadRequest,
				Body:   []byte("namespace cannot be empty"),
			})
		}

		// create metric names
		var metricNames *metricNamesResponse
		var err error
		if qm.RunnableQuery.ResourceName == "" && *qm.RunnableQuery.Namespace == "*" {
			metricNames, err = d.getNamespacesContainingResource(&qm, &sender)
			if err != nil {
				return err
			}
		} else if qm.RunnableQuery.ResourceName == "*" && *qm.RunnableQuery.Namespace != "*" && *qm.RunnableQuery.Namespace != "" {
			metricNames, err = d.getResourcesInNamespace(&qm, &sender)
			if err != nil {
				return err
			}
		} else {
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusBadRequest,
				Body:   []byte("query format is invalid"),
			})
		}

		// return metric names
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

func (d *Datasource) getNamespacesContainingResource(qm *models.QueryModel, sender *backend.CallResourceResponseSender) (*metricNamesResponse, error) {
	metricNames := &metricNamesResponse{}
	switch qm.RunnableQuery.ResourceType {
	case models.PipelineResourceType:
		namespacesWithPipelines, err := d.ListNamespacesWithPipelines(d.settings.Namespace)
		if err != nil {
			backend.Logger.Error("error listing namespaces with pipelines", "err", err)
			return nil, (*sender).Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
				Body:   []byte("error listing namespaces with pipelines"),
			})
		}
		metricNames.MetricNames = namespacesWithPipelines
	case models.VertexResourceType:
		namespacesWithVertices, err := d.ListNamespacesWithVertices(d.settings.Namespace)
		if err != nil {
			backend.Logger.Error("error listing namespaces with vertices", "err", err)
			return nil, (*sender).Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
				Body:   []byte("error listing namespaces with vertices"),
			})
		}
		metricNames.MetricNames = namespacesWithVertices
	case models.IsbsvcResourceType:
		namespacesWithInterStepBufferServices, err := d.ListNamespacesWithInterStepBufferServices(d.settings.Namespace)
		if err != nil {
			backend.Logger.Error("error listing namespaces with isbsvc", "err", err)
			return nil, (*sender).Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
				Body:   []byte("error listing namespaces with isbsvc"),
			})
		}
		metricNames.MetricNames = namespacesWithInterStepBufferServices
	default:
		backend.Logger.Error("unknown resource type, this shouldn't ever happen", "resourceType", qm.RunnableQuery.ResourceType)
		return nil, (*sender).Send(&backend.CallResourceResponse{
			Status: http.StatusInternalServerError,
			Body:   []byte("error listing namespaces, resource type unknown"),
		})
	}
	return metricNames, nil
}

func (d *Datasource) getResourcesInNamespace(qm *models.QueryModel, sender *backend.CallResourceResponseSender) (*metricNamesResponse, error) {
	metricNames := &metricNamesResponse{}
	switch qm.RunnableQuery.ResourceType {
	case models.PipelineResourceType:
		pipelinesInNamespace, err := d.ListPipelines(*qm.RunnableQuery.Namespace)
		if err != nil {
			backend.Logger.Error("error listing pipelines in namespace", "err", err)
			return nil, (*sender).Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
				Body:   []byte("error listing pipelines in namespace"),
			})
		}
		pipelineNamesInNamespace := make([]string, len(pipelinesInNamespace))
		for i := range pipelinesInNamespace {
			pipelineNamesInNamespace[i] = pipelinesInNamespace[i].Name
		}
		metricNames.MetricNames = pipelineNamesInNamespace
	case models.VertexResourceType:
		verticesInNamespace, err := d.ListPipelineVertices(*qm.RunnableQuery.Namespace, *qm.RunnableQuery.Pipeline)
		if err != nil {
			backend.Logger.Error("error listing pipeline vertices in namespace", "err", err)
			return nil, (*sender).Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
				Body:   []byte("error listing pipeline vertices in namespace"),
			})
		}
		vertexNamesInNamespace := make([]string, len(verticesInNamespace))
		for i := range verticesInNamespace {
			vertexNamesInNamespace[i] = verticesInNamespace[i].Labels[dfv1.KeyVertexName]
		}
		metricNames.MetricNames = vertexNamesInNamespace
	case models.IsbsvcResourceType:
		isbsvcsInNamespace, err := d.ListInterStepBufferServices(*qm.RunnableQuery.Namespace)
		if err != nil {
			backend.Logger.Error("error listing isbsvcs in namespace", "err", err)
			return nil, (*sender).Send(&backend.CallResourceResponse{
				Status: http.StatusInternalServerError,
				Body:   []byte("error listing isbsvcs in namespace"),
			})
		}
		isbsvcNamesInNamespace := make([]string, len(isbsvcsInNamespace))
		for i := range isbsvcsInNamespace {
			isbsvcNamesInNamespace[i] = isbsvcsInNamespace[i].Name
		}
		metricNames.MetricNames = isbsvcNamesInNamespace
	default:
		backend.Logger.Error("unknown resource type, this shouldn't ever happen", "resourceType", qm.RunnableQuery.ResourceType)
		return nil, (*sender).Send(&backend.CallResourceResponse{
			Status: http.StatusInternalServerError,
			Body:   []byte("error listing namespaces, resource type unknown"),
		})
	}
	return metricNames, nil
}
