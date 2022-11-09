package query

import (
	"context"
	"errors"
	"github.com/dseapy/grafana-numaflow-datasource/pkg/models"
	"github.com/dseapy/grafana-numaflow-datasource/pkg/query/scenario"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func RunQuery(_ context.Context, settings models.PluginSettings, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}
	var qm models.QueryModel
	backend.Logger.Debug("query json %v", string(query.JSON))
	if response.Error = qm.Unmarshall(query.JSON); response.Error != nil {
		return response
	}

	// validate
	if settings.Namespaced {
		qm.RunnableQuery.Namespace = &settings.Namespace
	}
	if qm.RunnableQuery.Namespace == nil {
		qm.RunnableQuery.Namespace = pointer.String(v1.NamespaceAll)
	}
	if *qm.RunnableQuery.Namespace == v1.NamespaceAll && qm.RunnableQuery.ResourceName != "*" {
		response.Error = errors.New(`"namespace" must be provided when requesting a single pipeline, vertex, or isbsvc by name`)
		return response
	}

	// create frames
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
