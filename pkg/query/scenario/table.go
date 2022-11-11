package scenario

import (
	"errors"
	"github.com/dseapy/grafana-numaflow-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	dfv1 "github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/strings/slices"
	"strings"
	"time"
)

func newTableFrames(nfClient *NFClients, query models.RunnableQuery) (data.Frames, error) {
	queryNamespace := getNamespace(&query)
	queryFilterNamespaces := getFilterNamespaces(&query)
	switch query.ResourceType {
	case models.PipelineResourceType:
		pipelines := []v1alpha1.Pipeline{}
		if *query.Pipeline == "*" {
			p, err := nfClient.ListPipelines(queryNamespace)
			if err != nil {
				return nil, err
			}
			for i := range p {
				if slices.Contains(queryFilterNamespaces, p[i].Namespace) {
					pipelines = append(pipelines, p[i])
				}
			}
		} else {
			pipeline, err := nfClient.GetPipeline(queryNamespace, *query.Pipeline)
			if err != nil {
				return nil, err
			}
			pipelines = []v1alpha1.Pipeline{*pipeline}
		}
		namespaces := make([]string, len(pipelines))
		names := make([]string, len(pipelines))
		phases := make([]string, len(pipelines))
		numVertices := make([]*uint32, len(pipelines))
		numSources := make([]*uint32, len(pipelines))
		numSinks := make([]*uint32, len(pipelines))
		numUDFs := make([]*uint32, len(pipelines))
		creationTime := make([]time.Time, len(pipelines))
		for i := range pipelines {
			namespaces[i] = pipelines[i].Namespace
			names[i] = pipelines[i].Name
			phases[i] = string(pipelines[i].Status.Phase)
			numVertices[i] = pipelines[i].Status.VertexCount
			numSources[i] = pipelines[i].Status.SourceCount
			numSinks[i] = pipelines[i].Status.SinkCount
			numUDFs[i] = pipelines[i].Status.UDFCount
			creationTime[i] = pipelines[i].CreationTimestamp.Time
		}

		fields := []*data.Field{
			data.NewField("namespace", nil, namespaces),
			data.NewField("name", nil, names),
			data.NewField("phase", nil, phases),
			data.NewField("vertices", nil, numVertices),
			data.NewField("sources", nil, numSources),
			data.NewField("sinks", nil, numSinks),
			data.NewField("UDFs", nil, numUDFs),
			data.NewField("creation time", nil, creationTime),
		}
		return data.Frames{data.NewFrame("pipelines", fields...)}, nil
	case models.VertexResourceType:
		vertices := []v1alpha1.Vertex{}
		if *query.Vertex == "*" {
			queryFilterPipelines := getFilterPipelines(&query)
			v, err := nfClient.ListVertices(queryNamespace)
			if err != nil {
				return nil, err
			}
			for i := range v {
				if slices.Contains(queryFilterNamespaces, v[i].Namespace) &&
					(*query.Pipeline == "" || slices.Contains(queryFilterPipelines, v[i].Labels[dfv1.KeyPipelineName])) {
					vertices = append(vertices, v[i])
				}
			}
		} else {
			vertex, err := nfClient.GetPipelineVertex(queryNamespace, *query.Pipeline, *query.Vertex)
			if err != nil {
				return nil, err
			}
			vertices = []v1alpha1.Vertex{*vertex}
		}
		namespaces := make([]string, len(vertices))
		names := make([]string, len(vertices))
		phases := make([]string, len(vertices))
		creationTime := make([]time.Time, len(vertices))
		for i := range vertices {
			namespaces[i] = vertices[i].Namespace
			names[i] = vertices[i].Labels[dfv1.KeyVertexName]
			phases[i] = string(vertices[i].Status.Phase)
			creationTime[i] = vertices[i].CreationTimestamp.Time
		}

		fields := []*data.Field{
			data.NewField("namespace", nil, namespaces),
			data.NewField("name", nil, names),
			data.NewField("phase", nil, phases),
			data.NewField("creation time", nil, creationTime),
		}
		return data.Frames{data.NewFrame("vertices", fields...)}, nil
	case models.IsbsvcResourceType:
		isbsvcs := []v1alpha1.InterStepBufferService{}
		if *query.InterStepBufferService == "*" {
			is, err := nfClient.ListInterStepBufferServices(queryNamespace)
			if err != nil {
				return nil, err
			}
			for i := range is {
				if slices.Contains(queryFilterNamespaces, is[i].Namespace) {
					isbsvcs = append(isbsvcs, is[i])
				}
			}
		} else {
			isbsvc, err := nfClient.GetInterStepBufferService(queryNamespace, *query.InterStepBufferService)
			if err != nil {
				return nil, err
			}
			isbsvcs = []v1alpha1.InterStepBufferService{*isbsvc}
		}
		namespaces := make([]string, len(isbsvcs))
		names := make([]string, len(isbsvcs))
		phases := make([]string, len(isbsvcs))
		creationTime := make([]time.Time, len(isbsvcs))
		for i := range isbsvcs {
			namespaces[i] = isbsvcs[i].Namespace
			names[i] = isbsvcs[i].Name
			phases[i] = string(isbsvcs[i].Status.Phase)
			creationTime[i] = isbsvcs[i].CreationTimestamp.Time
		}

		fields := []*data.Field{
			data.NewField("namespace", nil, namespaces),
			data.NewField("name", nil, names),
			data.NewField("phase", nil, phases),
			data.NewField("creation time", nil, creationTime),
		}
		return data.Frames{data.NewFrame("isbsvcs", fields...)}, nil
	}

	return nil, errors.New("unknown query resource type, this shouldn't happen")
}

func getNamespace(q *models.RunnableQuery) string {
	if isMultiNamespaceFilter(q) {
		return v1.NamespaceAll
	}
	return *q.Namespace
}

func isMultiNamespaceFilter(q *models.RunnableQuery) bool {
	return strings.Contains(*q.Namespace, ",")
}

func isMultiPipelineFilter(q *models.RunnableQuery) bool {
	return strings.Contains(*q.Pipeline, ",")
}

func getFilterNamespaces(q *models.RunnableQuery) []string {
	if !isMultiNamespaceFilter(q) {
		return []string{*q.Namespace}
	}
	ns := *q.Namespace
	ns = strings.ReplaceAll(ns, "{", "")
	ns = strings.ReplaceAll(ns, "}", "")
	return strings.Split(ns, ",")
}

func getFilterPipelines(q *models.RunnableQuery) []string {
	if !isMultiPipelineFilter(q) {
		return []string{*q.Pipeline}
	}
	pl := *q.Pipeline
	pl = strings.ReplaceAll(pl, "{", "")
	pl = strings.ReplaceAll(pl, "}", "")
	return strings.Split(pl, ",")
}
