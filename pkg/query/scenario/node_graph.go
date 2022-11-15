package scenario

import (
	"errors"
	"fmt"
	"github.com/dseapy/grafana-numaflow-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/numaproj/numaflow/pkg/apis/numaflow/v1alpha1"
	"github.com/numaproj/numaflow/pkg/isb"
	"math"
	"strconv"
	"time"
)

// Example on how you can structure data frames when returning node graph data.
func newNodeGraphFrames(nfClient *NFClients, query models.RunnableQuery) (data.Frames, error) {
	if query.ResourceType != models.PipelineResourceType {
		return nil, errors.New("node graph currently only supports pipelines")
	}
	if *query.Pipeline == "*" || query.IsMultiPipelineFilter() {
		return nil, errors.New("node graph currently only supports a single pipeline")
	}
	queryNamespace := query.GetNamespace()
	pipeline, err := nfClient.GetPipeline(queryNamespace, *query.Pipeline)
	if err != nil {
		return nil, err
	}
	if pipeline == nil {
		return nil, errors.New("retrieved pipeline was nil")
	}

	// get vertices & edges
	vertices, err := nfClient.ListPipelineVertices(queryNamespace, *query.Pipeline)
	if err != nil {
		return nil, err
	}
	edges, err := nfClient.ListPipelineEdges(queryNamespace, *query.Pipeline)
	if err != nil {
		return nil, err
	}

	// declare vertex and edge metrics
	vertexIDs := make([]string, len(vertices))
	vertexTitles := make([]string, len(vertices))
	vertexSubtitles := make([]string, len(vertices))       // # pods
	vertexMainStats := make([]string, len(vertices))       // processing rate
	vertexSecondaryStats := make([]*string, len(vertices)) // watermark
	vertexArcSuccess := make([]float32, len(vertices))
	vertexArcFailure := make([]float32, len(vertices))
	vertexArcNeutral := make([]float32, len(vertices))

	edgeIDs := make([]string, len(edges))
	edgeSources := make([]string, len(edges))
	edgeTargets := make([]string, len(edges))
	edgeMainStats := make([]string, len(edges))
	edgeSecondaryStats := make([]string, len(edges))

	// populate vertex and edge metrics
	for i := range vertices {
		vertexIDs[i] = vertices[i].Spec.Name
		vertexTitles[i] = vertices[i].Spec.Name
		specReplicas := int32(1)
		if vertices[i].Spec.Replicas != nil {
			specReplicas = *vertices[i].Spec.Replicas
		}
		vertexSubtitles[i] = fmt.Sprintf("%d/%d", vertices[i].Status.Replicas, specReplicas)
		vMetrics, err := nfClient.GetVertexMetrics(queryNamespace, vertices[i].Spec.PipelineName, vertices[i].Spec.Name)
		if err != nil {
			backend.Logger.Error("failed to retrieve metrics for vertex", "namespace", queryNamespace, "pipeline", vertices[i].Spec.PipelineName, "vertex", vertices[i].Spec.Name)
		} else {
			// Avg rate and pending for autoscaling are both in the map with key "default", see "pkg/metrics/metrics.go".
			rate, existing := vMetrics.ProcessingRates["default"]
			if !existing || rate < 0 || rate == isb.RateNotAvailable { // Rate not available
				backend.Logger.Debug("processing rate not available for vertex", "namespace", queryNamespace, "pipeline", vertices[i].Spec.PipelineName, "vertex", vertices[i].Spec.Name)
			} else {
				vertexMainStats[i] = strconv.FormatFloat(roundFloat(rate, 2), 'f', -1, 64) + " msg/s"
			}
		}
		vWatermark, err := nfClient.GetVertexWatermark(queryNamespace, vertices[i].Spec.PipelineName, vertices[i].Spec.Name)
		if err != nil {
			backend.Logger.Error("failed to retrieve watermark for vertex", "namespace", queryNamespace, "pipeline", vertices[i].Spec.PipelineName, "vertex", vertices[i].Spec.Name)
		} else {
			if vWatermark.Watermark != nil {
				t := time.UnixMilli(*vWatermark.Watermark).Format("2006-01-02T15:04:05.000Z")
				vertexSecondaryStats[i] = &t
			}
		}
		if vertices[i].Status.Phase == v1alpha1.VertexPhaseFailed {
			vertexArcSuccess[i] = float32(0.0)
			vertexArcFailure[i] = float32(1.0)
			vertexArcNeutral[i] = float32(0.0)
		} else if vertices[i].Status.Phase == v1alpha1.VertexPhaseSucceeded ||
			vertices[i].Status.Phase == v1alpha1.VertexPhaseRunning {
			vertexArcSuccess[i] = float32(1.0)
			vertexArcFailure[i] = float32(0.0)
			vertexArcNeutral[i] = float32(0.0)
		} else {
			vertexArcSuccess[i] = float32(0.0)
			vertexArcFailure[i] = float32(0.0)
			vertexArcNeutral[i] = float32(1.0)
		}
	}
	for i := range edges {
		if edges[i].FromVertex == nil {
			return nil, errors.New("edge from-vertex was nil")
		}
		edgeSources[i] = *edges[i].FromVertex
		if edges[i].ToVertex == nil {
			return nil, errors.New("edge to-vertex was nil")
		}
		edgeTargets[i] = *edges[i].ToVertex
		edgeIDs[i] = fmt.Sprintf("%s-%s", edgeSources[i], edgeTargets[i])
		if edges[i].PendingCount != nil && edges[i].AckPendingCount != nil {
			edgeMainStats[i] = fmt.Sprintf("%d", *edges[i].PendingCount+*edges[i].AckPendingCount)
		}
		if edges[i].BufferUsage != nil {
			edgeSecondaryStats[i] = strconv.FormatFloat(*edges[i].BufferUsage, 'f', -1, 64) + "%"
		}
	}

	// set return fields with vertex and edge metrics
	arcSuccessConfig := &data.FieldConfig{Color: map[string]interface{}{"mode": "fixed", "fixedColor": "green"}}
	arcFailureConfig := &data.FieldConfig{Color: map[string]interface{}{"mode": "fixed", "fixedColor": "red"}}
	vertexFields := []*data.Field{
		data.NewField("id", nil, vertexIDs),
		data.NewField("title", nil, vertexTitles),
		data.NewField("subtitle", nil, vertexSubtitles),
		data.NewField("mainstat", nil, vertexMainStats),
		data.NewField("secondarystat", nil, vertexSecondaryStats),
		data.NewField("arc__success", nil, vertexArcSuccess).SetConfig(arcSuccessConfig),
		data.NewField("arc__failure", nil, vertexArcFailure).SetConfig(arcFailureConfig),
		data.NewField("arc__neutral", nil, vertexArcNeutral),
	}
	verticesFrame := data.NewFrame("nodes", vertexFields...)

	edgeFields := []*data.Field{
		data.NewField("id", nil, edgeIDs),
		data.NewField("source", nil, edgeSources),
		data.NewField("target", nil, edgeTargets),
		data.NewField("mainstat", nil, edgeMainStats),
		data.NewField("secondarystat", nil, edgeSecondaryStats),
	}
	edgesFrame := data.NewFrame("edges", edgeFields...)

	return data.Frames{verticesFrame, edgesFrame}, nil
}

// TODO: move to a utils
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
