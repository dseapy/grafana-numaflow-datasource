package query

import (
	"context"
	"encoding/json"

	"github.com/dseapy/grafana-numaflow-datasource/pkg/models"
	"github.com/dseapy/grafana-numaflow-datasource/pkg/query/scenario"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func RunQuery(_ context.Context, settings models.PluginSettings, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	// Unmarshal the JSON into our queryModel.
	var qm models.QueryModel

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	// We are not using the RunnableQuery in this example because we are generating
	// static data depending on the query type. We still want to show case how to
	// support macros/server side variables in your queries.
	frames := scenario.NewDataFrames(query)
	if len(frames) == 0 {
		return response
	}

	for _, frame := range frames {
		// Assign the refId from the query to the reply data frame to make it
		// easier to track.
		frame.RefID = query.RefID
	}

	// Add the frames to the response.
	response.Frames = append(response.Frames, frames...)

	return response
}
