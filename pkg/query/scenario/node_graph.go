package scenario

import (
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Example on how you can structure data frames when returning node graph data.
func newNodeGraphFrames() data.Frames {
	nodeFields := []*data.Field{
		// required
		data.NewField("id", nil, []string{"my-node-id1", "my-node-id2", "my-node-id3"}),
		// optional
		data.NewField("title", nil, []string{"my-node-title1", "my-node-title2", ""}),
		data.NewField("subtitle", nil, []string{"my-node-subtitle1", "my-node-subtitle2", ""}),
		data.NewField("mainstat", nil, []string{"my-node-mainstat1", "my-node-mainstat2", ""}),
		data.NewField("secondarystat", nil, []string{"my-node-secondarystat1", "my-node-secondarystat2", ""}),
		data.NewField("arc__foo", nil, []float32{0.3, 0.5, 1.0}),
		data.NewField("arc__bar", nil, []float32{0.7, 0.5, 1.0}),
		data.NewField("detail__zed", nil, []string{"my-node-detail-zed1", "my-node-detail-zed2", ""}),
	}
	nodesFrame := data.NewFrame("nodes", nodeFields...)

	edgeFields := []*data.Field{
		// required
		data.NewField("id", nil, []string{"my-edge-id1", "my-edge-id2"}),
		data.NewField("source", nil, []string{"my-node-id1", "my-node-id2"}),
		data.NewField("target", nil, []string{"my-node-id2", "my-node-id3"}),
		// optional
		data.NewField("mainstat", nil, []string{"my-edge-mainstat1", "my-edge-mainstat2"}),
		data.NewField("secondarystat", nil, []string{"my-edge-secondarystat1", "my-edge-secondarystat2"}),
		data.NewField("detail__zed", nil, []string{"my-edge-detail-zed1", "my-edge-detail-zed2"}),
	}
	edgesFrame := data.NewFrame("edges", edgeFields...)

	return data.Frames{nodesFrame, edgesFrame}
}
