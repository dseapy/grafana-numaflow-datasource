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

if "vertex" exists, "pipeline" must exist.

"*": give me all of those
"my-foo": give me that one
"": ignore, but relevant

DATA_QUERY
----------
For all pipelines all namespaces: {"namespace":"","pipeline":"*"}
For all vertices all pipelines: {"namespace":"","pipeline":"","vertex":"*"}
For all isbsvcs all namespace: {"namespace":"","isbsvc":"*"}
For all pipelines in namespace: {"namespace":"$namespace","pipeline":"*"}
For all vertices in namespace: {"namespace":"$namespace","pipeline":"","vertex":"*"}
For all vertices in pipeline: {"namespace":"$namespace","pipeline":"my-pl","vertex":"*"}
For all isbsvcs in namespace: {"namespace":"$namespace","isbsvc":"*"}
For a single pipeline: {"namespace":"$namespace","pipeline":"my-pl"}
For a single vertex: {"namespace":"$namespace","pipeline":"my-pl","vertex":"my-vertex"}
For a single isbsvc: {"namespace":"$namespace","isbsvc":"my-isbsvc"}

METRIC_NAME_QUERY
-----------------
For all pipelines in namespace: {"namespace":"$namespace","pipeline":"*"}
For all vertices in namespace: {"namespace":"$namespace","pipeline":"","vertex":"*"}
For all vertices in pipeline: {"namespace":"$namespace","pipeline":"$pipeline","vertex":"*"}
For all isbsvcs in namespace: {"namespace":"$namespace","isbsvc":"*"}
For all namespaces containing pipelines: {"namespace":"*","pipeline":""}
For all namespaces containing vertices: {"namespace":"*","pipeline":"","vertex":""}
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

	if qm.RunnableQuery.Pipeline != nil {
		if qm.RunnableQuery.Vertex != nil {
			qm.RunnableQuery.ResourceName = *qm.RunnableQuery.Vertex
			qm.RunnableQuery.ResourceType = VertexResourceType
		} else {
			qm.RunnableQuery.ResourceName = *qm.RunnableQuery.Pipeline
			qm.RunnableQuery.ResourceType = PipelineResourceType
		}
	} else if qm.RunnableQuery.InterStepBufferService != nil {
		qm.RunnableQuery.ResourceName = *qm.RunnableQuery.InterStepBufferService
		qm.RunnableQuery.ResourceType = IsbsvcResourceType
	} else {
		// TODO: better error reporting
		return errors.New("cannot unmarshal json")
	}

	return nil
}
