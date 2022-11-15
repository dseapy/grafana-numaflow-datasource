package scenario

import (
	"errors"
	"github.com/dseapy/numaflow-datasource/pkg/client"
	"github.com/dseapy/numaflow-datasource/pkg/query"
	"k8s.io/utils/pointer"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	"github.com/numaproj/numaflow/pkg/isb"
	"k8s.io/utils/strings/slices"
)

func newTableFrames(client *client.Client, rq query.RunnableQuery) (data.Frames, error) {
	switch rq.ResourceType {
	case query.PipelineResourceType:
		return newPipelineTableFrames(client, rq)
	case query.VertexResourceType:
		return newVertexTableFrames(client, rq)
	case query.IsbsvcResourceType:
		return newIsbsvcTableFrames(client, rq)
	}
	return nil, errors.New("unknown query resource type, this shouldn't happen")
}

func newPipelineTableFrames(nfClient *client.Client, rq query.RunnableQuery) (data.Frames, error) {
	queryNamespace := rq.GetNamespace()
	queryFilterNamespaces := rq.GetFilterNamespaces()
	pipelines := []v1alpha1.Pipeline{}
	if *rq.Pipeline == "*" {
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
		pipeline, err := nfClient.GetPipeline(queryNamespace, *rq.Pipeline)
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
}

func newVertexTableFrames(nfClient *client.Client, rq query.RunnableQuery) (data.Frames, error) {
	queryNamespace := rq.GetNamespace()
	queryFilterNamespaces := rq.GetFilterNamespaces()
	vertices := []v1alpha1.Vertex{}
	queryPipeline := *rq.Pipeline
	if *rq.Vertex == "*" {
		queryFilterPipelines := rq.GetFilterPipelines()
		v, err := nfClient.ListVertices(queryNamespace)
		if err != nil {
			return nil, err
		}
		for i := range v {
			if slices.Contains(queryFilterNamespaces, v[i].Namespace) &&
				(queryPipeline == "" || slices.Contains(queryFilterPipelines, v[i].Spec.PipelineName)) {
				vertices = append(vertices, v[i])
			}
		}
	} else {
		vertex, err := nfClient.GetPipelineVertex(queryNamespace, queryPipeline, *rq.Vertex)
		if err != nil {
			return nil, err
		}
		vertices = []v1alpha1.Vertex{*vertex}
	}
	namespaces := make([]string, len(vertices))
	pipelineNames := make([]string, len(vertices))
	names := make([]string, len(vertices))
	vtype := make([]string, len(vertices))
	desiredReplicas := make([]*int32, len(vertices))
	replicas := make([]uint32, len(vertices))
	phases := make([]string, len(vertices))
	processingRate := make([]*float64, len(vertices))
	pendingMessages := make([]*int64, len(vertices))
	watermark := make([]*time.Time, len(vertices))
	cpuUsage := make([]*int64, len(vertices))
	memoryUsage := make([]*int64, len(vertices))
	creationTime := make([]time.Time, len(vertices))
	for i := range vertices {
		namespaces[i] = vertices[i].Namespace
		pipelineNames[i] = vertices[i].Spec.PipelineName
		names[i] = vertices[i].Spec.Name
		if vertices[i].IsASource() {
			vtype[i] = "source"
		} else if vertices[i].IsASink() {
			vtype[i] = "sink"
		} else if vertices[i].IsMapUDF() {
			vtype[i] = "udf (map)"
		} else if vertices[i].IsReduceUDF() {
			vtype[i] = "udf (reduce)"
		}
		phases[i] = string(vertices[i].Status.Phase)
		desiredReplicas[i] = vertices[i].Spec.Replicas
		replicas[i] = vertices[i].Status.Replicas
		vMetrics, err := nfClient.GetVertexMetrics(queryNamespace, queryPipeline, vertices[i].Spec.Name)
		if err != nil {
			backend.Logger.Error("failed to retrieve metrics for vertex", "namespace", queryNamespace, "pipeline", queryPipeline, "vertex", vertices[i].Spec.Name)
		} else {
			// Avg rate and pending for autoscaling are both in the map with key "default", see "pkg/metrics/metrics.go".
			rate, existing := vMetrics.ProcessingRates["default"]
			if !existing || rate < 0 || rate == isb.RateNotAvailable { // Rate not available
				backend.Logger.Debug("processing rate not available for vertex", "namespace", queryNamespace, "pipeline", queryPipeline, "vertex", vertices[i].Spec.Name)
			} else {
				processingRate[i] = pointer.Float64(rate)
			}
			pending, existing := vMetrics.Pendings["default"]
			if !existing || pending < 0 || pending == isb.PendingNotAvailable {
				backend.Logger.Debug("pending not available for vertex", "namespace", queryNamespace, "pipeline", queryPipeline, "vertex", vertices[i].Spec.Name)
			} else {
				pendingMessages[i] = pointer.Int64(pending)
			}
		}
		vWatermark, err := nfClient.GetVertexWatermark(queryNamespace, queryPipeline, vertices[i].Spec.Name)
		if err != nil {
			backend.Logger.Error("failed to retrieve watermark for vertex", "namespace", queryNamespace, "pipeline", queryPipeline, "vertex", vertices[i].Spec.Name)
		} else {
			if vWatermark.Watermark != nil {
				t := time.UnixMilli(*vWatermark.Watermark)
				watermark[i] = &t
			}
		}
		pods, err := nfClient.ListVertexPods(queryNamespace, queryPipeline, vertices[i].Spec.Name)
		if err != nil {
			backend.Logger.Error("failed to retrieve pods for vertex", "namespace", queryNamespace, "pipeline", queryPipeline, "vertex", vertices[i].Spec.Name)
		} else {
			pCpu := int64(0)
			pMemory := int64(0)
			for pi := range pods {
				pMetrics, err := nfClient.GetPodMetrics(queryNamespace, pods[pi].Name)
				if err != nil {
					backend.Logger.Error("failed to retrieve pod metrics for pod", "namespace", queryNamespace, "pod", pods[pi].Name)
					pCpu = int64(0)
					pMemory = int64(0)
					break
				}
				for _, c := range pMetrics.Containers {
					pCpu += c.Usage.Cpu().MilliValue()
					pMemory += c.Usage.Memory().ScaledValue(6)
				}
			}
			if pCpu != int64(0) && pMemory != int64(0) {
				cpuUsage[i] = &pCpu
				memoryUsage[i] = &pMemory
			}
		}
		creationTime[i] = vertices[i].CreationTimestamp.Time
	}

	fields := []*data.Field{
		data.NewField("namespace", nil, namespaces),
		data.NewField("pipeline", nil, pipelineNames),
		data.NewField("name", nil, names),
		data.NewField("type", nil, vtype),
		data.NewField("phase", nil, phases),
		data.NewField("replicas", nil, replicas),
		data.NewField("desired replicas", nil, desiredReplicas),
		data.NewField("processing rate", nil, processingRate),
		data.NewField("pending messages", nil, pendingMessages),
		data.NewField("watermark", nil, watermark),
		data.NewField("cpu usage", nil, cpuUsage),
		data.NewField("memory usage", nil, memoryUsage),
		data.NewField("creation time", nil, creationTime),
	}
	return data.Frames{data.NewFrame("vertices", fields...)}, nil
}

func newIsbsvcTableFrames(nfClient *client.Client, rq query.RunnableQuery) (data.Frames, error) {
	queryNamespace := rq.GetNamespace()
	queryFilterNamespaces := rq.GetFilterNamespaces()
	isbsvcs := []v1alpha1.InterStepBufferService{}
	if *rq.InterStepBufferService == "*" {
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
		isbsvc, err := nfClient.GetInterStepBufferService(queryNamespace, *rq.InterStepBufferService)
		if err != nil {
			return nil, err
		}
		isbsvcs = []v1alpha1.InterStepBufferService{*isbsvc}
	}
	namespaces := make([]string, len(isbsvcs))
	names := make([]string, len(isbsvcs))
	itype := make([]string, len(isbsvcs))
	phases := make([]string, len(isbsvcs))
	creationTime := make([]time.Time, len(isbsvcs))
	for i := range isbsvcs {
		namespaces[i] = isbsvcs[i].Namespace
		names[i] = isbsvcs[i].Name
		if isbsvcs[i].Spec.Redis != nil {
			if isbsvcs[i].Spec.Redis.External != nil {
				itype[i] = "redis (external)"
			} else {
				itype[i] = "redis (internal)"
			}
		} else if isbsvcs[i].Spec.JetStream != nil {
			itype[i] = "jetstream"
		}
		phases[i] = string(isbsvcs[i].Status.Phase)
		creationTime[i] = isbsvcs[i].CreationTimestamp.Time
	}

	fields := []*data.Field{
		data.NewField("namespace", nil, namespaces),
		data.NewField("name", nil, names),
		data.NewField("type", nil, itype),
		data.NewField("phase", nil, phases),
		data.NewField("creation time", nil, creationTime),
	}
	return data.Frames{data.NewFrame("isbsvcs", fields...)}, nil
}
