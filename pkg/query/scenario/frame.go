package scenario

import (
	"github.com/dseapy/grafana-numaflow-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func NewDataFrames(query backend.DataQuery, runnableQuery models.RunnableQuery) data.Frames {
	switch query.QueryType {
	case Table:
		return newTableFrames(runnableQuery)
	case NodeGraph:
		return newNodeGraphFrames(runnableQuery)
	}

	return nil
}
