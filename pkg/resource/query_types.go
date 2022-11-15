package resource

import (
	"encoding/json"
	"net/http"
)

const (
	TableQueryType     string = "Table"
	NodeGraphQueryType string = "NodeGraph"

	QueryTypesAPIPath   = "/query-types"
	QueryTypesAPIMethod = http.MethodGet
)

type queryTypes struct {
	QueryTypes []string `json:"queryTypes"`
}

func QueryTypes() []string {
	return []string{
		TableQueryType,
		NodeGraphQueryType,
	}
}

func QueryTypesJson() ([]byte, error) {
	q := &queryTypes{
		QueryTypes: QueryTypes(),
	}
	j, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}
	return j, nil
}
