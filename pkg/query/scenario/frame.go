package scenario

import (
	"errors"
	"github.com/dseapy/grafana-numaflow-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func NewDataFrames(nfClient *NFClients, query backend.DataQuery, runnableQuery models.RunnableQuery) (data.Frames, error) {
	switch query.QueryType {
	case Table:
		return newTableFrames(nfClient, runnableQuery)
	case NodeGraph:
		return newNodeGraphFrames(nfClient, runnableQuery)
	}

	return nil, errors.New("unsupported query type")
}
