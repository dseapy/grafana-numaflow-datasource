package scenario

import (
	"errors"
	"github.com/dseapy/numaflow-datasource/pkg/client"
	"github.com/dseapy/numaflow-datasource/pkg/query"
	"github.com/dseapy/numaflow-datasource/pkg/resource"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func NewDataFrames(nfClient *client.Client, query backend.DataQuery, runnableQuery query.RunnableQuery) (data.Frames, error) {
	switch query.QueryType {
	case resource.TableQueryType:
		return newTableFrames(nfClient, runnableQuery)
	case resource.NodeGraphQueryType:
		return newNodeGraphFrames(nfClient, runnableQuery)
	}

	return nil, errors.New("unsupported query type")
}
