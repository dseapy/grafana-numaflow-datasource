package models

import (
	"encoding/json"
	"errors"
	v1 "k8s.io/api/core/v1"
)

type QueryModel struct {
	RawQuery      string        `json:"rawQuery"`
	RunnableQuery RunnableQuery `json:"-"`
}

const (
	// ResourceAll is an argument used to request all resources of a particular resource type
	ResourceAll string = "*"
)

/*
RunnableQuery describes what data should be returned by the backend.

{"namespace":"my-ns"} may be optionally included to restrict namespace.
if "namespaced" datasource setting is true, the "namespace" datasource setting is always used istead.
if "namespaced" datasource setting is false, "namespace" is required in the query only when specifying a single pipeline/vertex/isbsvc.

Only one resource may be specified ("pipeline", "vertex", or "isbsvc")

For all pipelines: {"pipeline":"*"}
For a single pipeline: {"pipeline":"my-pl"}
For all vertices: {"vertex":"*"}
For a single vertex: {"vertex":"my-vertex"}
For all isbsvc: {"isbsvc":"*"}
For a single isbsvc: {"isbsvc":"my-isbsvc"}
*/
type RunnableQuery struct {
	Namespace              string `json:"namespace,omitempty"`
	Pipeline               string `json:"pipeline,omitempty"`
	Vertex                 string `json:"vertex,omitempty"`
	InterStepBufferService string `json:"isbsvc,omitempty"`
}

func (qm *QueryModel) Unmarshall(b []byte, settings PluginSettings) error {
	if err := json.Unmarshal(b, &qm); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(qm.RawQuery), &qm.RunnableQuery); err != nil {
		return err
	}
	if settings.Namespaced {
		qm.RunnableQuery.Namespace = settings.Namespace
	}

	numResourcesSpecified := 0
	resourceName := ""
	if qm.RunnableQuery.Pipeline != "" {
		resourceName = qm.RunnableQuery.Pipeline
		numResourcesSpecified++
	}
	if qm.RunnableQuery.Vertex != "" {
		resourceName = qm.RunnableQuery.Vertex
		numResourcesSpecified++
	}
	if qm.RunnableQuery.InterStepBufferService != "" {
		resourceName = qm.RunnableQuery.InterStepBufferService
		numResourcesSpecified++
	}
	if numResourcesSpecified != 1 {
		return errors.New(`must specify exactly one of the following in query: "pipeline", "vertex", "isbsvc"`)
	}
	if qm.RunnableQuery.Namespace == v1.NamespaceAll && resourceName != ResourceAll {
		return errors.New(`"namespace" must be provided when requesting a single pipeline, vertex, or isbsvc by name`)
	}

	return nil
}
