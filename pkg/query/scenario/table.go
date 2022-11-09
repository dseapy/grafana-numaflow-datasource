package scenario

import (
	"github.com/dseapy/grafana-numaflow-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Example on how you can structure your data frames when returning
// table based data.
func newTableFrames(query models.RunnableQuery) data.Frames {
	tempInside := []int64{25, 22, 19, 23, 22, 22, 18, 26, 24, 20}
	tempOutside := []int64{10, 8, 12, 9, 10, 11, 10, 9, 10, 9}

	fields := []*data.Field{
		data.NewField("temperature", data.Labels{"sensor": "outside"}, tempOutside),
		data.NewField("temperature", data.Labels{"sensor": "inside"}, tempInside),
	}

	return data.Frames{data.NewFrame("temperatures", fields...)}
}
