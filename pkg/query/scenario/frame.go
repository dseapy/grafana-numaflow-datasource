package scenario

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func NewDataFrames(query backend.DataQuery) data.Frames {
	switch query.QueryType {
	case TimeSeries:
		return newTimeSeriesFrames(query)
	case Table:
		return newTableFrames(query)
	case NodeGraph:
		return newNodeGraphFrames()
	}

	return nil
}
