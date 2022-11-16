package query

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

// RunnableQuery describes what data should be returned by the backend.
type RunnableQuery struct {
	Namespace              *string      `json:"namespace,omitempty"`
	Pipeline               *string      `json:"pipeline,omitempty"`
	Vertex                 *string      `json:"vertex,omitempty"`
	Pod                    *string      `json:"pod,omitempty"`
	InterStepBufferService *string      `json:"isbsvc,omitempty"`
	ResourceType           ResourceType `json:"-"`
	ResourceName           string       `json:"-"`
}

func (q *RunnableQuery) GetNamespace() string {
	if q.IsMultiNamespaceFilter() {
		return v1.NamespaceAll
	}
	return *q.Namespace
}

func (q *RunnableQuery) IsMultiNamespaceFilter() bool {
	return strings.Contains(*q.Namespace, ",")
}

func (q *RunnableQuery) IsMultiPipelineFilter() bool {
	return strings.Contains(*q.Pipeline, ",")
}

func (q *RunnableQuery) GetFilterNamespaces() []string {
	if !q.IsMultiNamespaceFilter() {
		return []string{*q.Namespace}
	}
	ns := *q.Namespace
	ns = strings.ReplaceAll(ns, "{", "")
	ns = strings.ReplaceAll(ns, "}", "")
	return strings.Split(ns, ",")
}

func (q *RunnableQuery) GetFilterPipelines() []string {
	if !q.IsMultiPipelineFilter() {
		return []string{*q.Pipeline}
	}
	pl := *q.Pipeline
	pl = strings.ReplaceAll(pl, "{", "")
	pl = strings.ReplaceAll(pl, "}", "")
	return strings.Split(pl, ",")
}
