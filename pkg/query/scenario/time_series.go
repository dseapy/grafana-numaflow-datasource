package scenario

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Example on how you can structure data frames when returning time
// series data.
func newTimeSeriesFrames(query backend.DataQuery) data.Frames {
	temperatures := []int64{25, 22, 19, 23, 22, 22, 18, 26, 24, 20}
	timestamps := timeStampsBetween(query.TimeRange, len(temperatures))

	fields := []*data.Field{
		data.NewField("time", nil, timestamps),
		data.NewField("values", nil, temperatures),
	}

	return data.Frames{data.NewFrame("temperatures", fields...)}
}
