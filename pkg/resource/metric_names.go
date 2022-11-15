package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dseapy/numaflow-datasource/pkg/client"
	"github.com/dseapy/numaflow-datasource/pkg/query"
	dfv1 "github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"

	v1 "k8s.io/api/core/v1"
)

const (
	MetricNamesAPIPath   = "/metric-names"
	MetricNamesAPIMethod = http.MethodPost
)

type metricNames struct {
	MetricNames []string `json:"metricNames"`
}

func MetricNamesJson(q *query.Query, c *client.Client) ([]byte, error) {
	// validate
	if q.RunnableQuery.Namespace == nil || *q.RunnableQuery.Namespace == v1.NamespaceAll {
		return nil, errors.New("namespace cannot be empty")
	}

	// create metric names
	var mn *metricNames
	var err error
	if q.RunnableQuery.ResourceName == "" && *q.RunnableQuery.Namespace == "*" {
		mn, err = getNamespacesContainingResource(q, c)
		if err != nil {
			return nil, err
		}
	} else if q.RunnableQuery.ResourceName == "*" && *q.RunnableQuery.Namespace != "*" && *q.RunnableQuery.Namespace != "" {
		mn, err = getResourcesInNamespace(q, c)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("query format is invalid")
	}

	// return metric names
	j, err := json.Marshal(mn)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func getNamespacesContainingResource(q *query.Query, c *client.Client) (*metricNames, error) {
	mn := &metricNames{}
	switch q.RunnableQuery.ResourceType {
	case query.PipelineResourceType:
		namespacesWithPipelines, err := c.ListNamespacesWithPipelines()
		if err != nil {
			return nil, err
		}
		mn.MetricNames = namespacesWithPipelines
	case query.VertexResourceType:
		namespacesWithVertices, err := c.ListNamespacesWithVertices()
		if err != nil {
			return nil, err
		}
		mn.MetricNames = namespacesWithVertices
	case query.IsbsvcResourceType:
		namespacesWithInterStepBufferServices, err := c.ListNamespacesWithInterStepBufferServices()
		if err != nil {
			return nil, err
		}
		mn.MetricNames = namespacesWithInterStepBufferServices
	default:
		return nil, errors.New(fmt.Sprintf("error listing namespaces, resource type unknown, %v", q.RunnableQuery.ResourceType))
	}
	return mn, nil
}

func getResourcesInNamespace(q *query.Query, c *client.Client) (*metricNames, error) {
	mn := &metricNames{}
	switch q.RunnableQuery.ResourceType {
	case query.PipelineResourceType:
		pipelinesInNamespace, err := c.ListPipelines(*q.RunnableQuery.Namespace)
		if err != nil {
			return nil, err
		}
		pipelineNamesInNamespace := make([]string, len(pipelinesInNamespace))
		for i := range pipelinesInNamespace {
			pipelineNamesInNamespace[i] = pipelinesInNamespace[i].Name
		}
		mn.MetricNames = pipelineNamesInNamespace
	case query.VertexResourceType:
		verticesInPipeline, err := c.ListPipelineVertices(*q.RunnableQuery.Namespace, *q.RunnableQuery.Pipeline)
		if err != nil {
			return nil, err
		}
		vertexNamesInPipeline := make([]string, len(verticesInPipeline))
		for i := range verticesInPipeline {
			vertexNamesInPipeline[i] = verticesInPipeline[i].Labels[dfv1.KeyVertexName]
		}
		mn.MetricNames = vertexNamesInPipeline
	case query.IsbsvcResourceType:
		isbsvcsInNamespace, err := c.ListInterStepBufferServices(*q.RunnableQuery.Namespace)
		if err != nil {
			return nil, err
		}
		isbsvcNamesInNamespace := make([]string, len(isbsvcsInNamespace))
		for i := range isbsvcsInNamespace {
			isbsvcNamesInNamespace[i] = isbsvcsInNamespace[i].Name
		}
		mn.MetricNames = isbsvcNamesInNamespace
	case query.PodResourceType:
		if q.RunnableQuery.Vertex != nil {
			podsInIsbsvc, err := c.ListVertexPods(*q.RunnableQuery.Namespace, *q.RunnableQuery.Pipeline, *q.RunnableQuery.Vertex)
			if err != nil {
				return nil, err
			}
			podNamesInIsbsvc := make([]string, len(podsInIsbsvc))
			for i := range podsInIsbsvc {
				podNamesInIsbsvc[i] = podsInIsbsvc[i].Name
			}
			mn.MetricNames = podNamesInIsbsvc
		} else if q.RunnableQuery.InterStepBufferService != nil {

		} else {
			return nil, errors.New(fmt.Sprintf("vertex or isbsvc must be provided when requesting pods"))
		}
	default:
		return nil, errors.New(fmt.Sprintf("error listing namespaces, resource type unknown, %v", q.RunnableQuery.ResourceType))
	}
	return mn, nil
}
