package models

import (
	"encoding/json"
	"errors"
)

type QueryModel struct {
	RawQuery      string        `json:"rawQuery"`
	RunnableQuery RunnableQuery `json:"-"`
}

type NumaflowResourceType string

const (
	NamespaceResourceType NumaflowResourceType = "namespace"
	PipelineResourceType  NumaflowResourceType = "pipeline"
	VertexResourceType    NumaflowResourceType = "vertex"
	IsbsvcResourceType    NumaflowResourceType = "isbsvc"
)

/*
RunnableQuery describes what data should be returned by the backend.

Only one non-namespace resource may be specified ("pipeline", "vertex", or "isbsvc")

DATA_QUERY
----------
For all pipelines all namespaces: {"namespace":"","pipeline":"*"}
For all vertices all namespace: {"namespace":"","vertex":"*"}
For all isbsvcs all namespace: {"namespace":"","isbsvc":"*"}
For all pipelines in namespace: {"namespace":"my-ns","pipeline":"*"}
For all vertices in namespace: {"namespace":"my-ns","vertex":"*"}
For all isbsvcs in namespace: {"namespace":"my-ns","isbsvc":"*"}
For a single pipeline: {"namespace":"my-ns","pipeline":"my-pl"}
For a single vertex: {"namespace":"my-ns","vertex":"my-vertex"}
For a single isbsvc: {"namespace":"my-ns","isbsvc":"my-isbsvc"}

METRIC_NAME_QUERY
-----------------
For all pipelines in namespace: {"namespace":"my-ns","pipeline":"*"}
For all vertices in namespace: {"namespace":"my-ns","vertex":"*"}
For all isbsvcs in namespace: {"namespace":"my-ns","isbsvc":"*"}
For all namespaces containing pipelines: {"namespace":"*","pipeline":""}
For all namespaces containing vertices: {"namespace":"*","vertex":""}
For all namespaces containing isbsvcs: {"namespace":"*","isbsvc":""}
*/
type RunnableQuery struct {
	Namespace              *string              `json:"namespace,omitempty"`
	Pipeline               *string              `json:"pipeline,omitempty"`
	Vertex                 *string              `json:"vertex,omitempty"`
	InterStepBufferService *string              `json:"isbsvc,omitempty"`
	ResourceType           NumaflowResourceType `json:"-"`
	ResourceName           string               `json:"-"`
}

func (qm *QueryModel) Unmarshall(b []byte) error {
	if err := json.Unmarshal(b, &qm); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(qm.RawQuery), &qm.RunnableQuery); err != nil {
		return err
	}

	numResourcesSpecified := 0
	if qm.RunnableQuery.Pipeline != nil {
		qm.RunnableQuery.ResourceName = *qm.RunnableQuery.Pipeline
		qm.RunnableQuery.ResourceType = PipelineResourceType
		numResourcesSpecified++
	}
	if qm.RunnableQuery.Vertex != nil {
		qm.RunnableQuery.ResourceName = *qm.RunnableQuery.Vertex
		qm.RunnableQuery.ResourceType = VertexResourceType
		numResourcesSpecified++
	}
	if qm.RunnableQuery.InterStepBufferService != nil {
		qm.RunnableQuery.ResourceName = *qm.RunnableQuery.InterStepBufferService
		qm.RunnableQuery.ResourceType = IsbsvcResourceType
		numResourcesSpecified++
	}
	if numResourcesSpecified != 1 {
		return errors.New(`must specify exactly one of the following in query: "pipeline", "vertex", "isbsvc"`)
	}

	return nil
}
