package query

import (
	"encoding/json"
	"errors"
)

type Query struct {
	RawQuery      string        `json:"rawQuery"`
	RunnableQuery RunnableQuery `json:"-"`
}

type ResourceType string

const (
	PipelineResourceType ResourceType = "pipeline"
	VertexResourceType   ResourceType = "vertex"
	IsbsvcResourceType   ResourceType = "isbsvc"
)

func (q *Query) Unmarshall(b []byte) error {
	if err := json.Unmarshal(b, &q); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(q.RawQuery), &q.RunnableQuery); err != nil {
		return err
	}

	if q.RunnableQuery.Pipeline != nil {
		if q.RunnableQuery.Vertex != nil {
			q.RunnableQuery.ResourceName = *q.RunnableQuery.Vertex
			q.RunnableQuery.ResourceType = VertexResourceType
		} else {
			q.RunnableQuery.ResourceName = *q.RunnableQuery.Pipeline
			q.RunnableQuery.ResourceType = PipelineResourceType
		}
	} else if q.RunnableQuery.InterStepBufferService != nil {
		q.RunnableQuery.ResourceName = *q.RunnableQuery.InterStepBufferService
		q.RunnableQuery.ResourceType = IsbsvcResourceType
	} else {
		// TODO: better error reporting
		return errors.New("cannot unmarshal json")
	}

	return nil
}
