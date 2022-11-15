package query

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

/*
RunnableQuery describes what data should be returned by the backend.

if "vertex" exists, "pipeline" must exist.

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
	Namespace              *string      `json:"namespace,omitempty"`
	Pipeline               *string      `json:"pipeline,omitempty"`
	Vertex                 *string      `json:"vertex,omitempty"`
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
