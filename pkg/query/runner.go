package query

import (
	"context"
	"github.com/dseapy/grafana-numaflow-datasource/pkg/models"
	"github.com/dseapy/grafana-numaflow-datasource/pkg/query/scenario"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func RunQuery(_ context.Context, settings models.PluginSettings, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}
	var qm models.QueryModel
	backend.Logger.Debug("query json %v", string(query.JSON))
	if response.Error = qm.Unmarshall(query.JSON, settings); response.Error != nil {
		return response
	}

	frames := scenario.NewDataFrames(query, qm.RunnableQuery)
	if len(frames) == 0 {
		return response
	}
	for _, frame := range frames {
		frame.RefID = query.RefID
	}
	response.Frames = append(response.Frames, frames...)
	return response
}
