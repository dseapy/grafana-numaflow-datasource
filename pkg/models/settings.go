package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type PluginSettings struct {
	DefaultNamespaced bool   `json:"namespaced"`
	DefaultNamespace  string `json:"namespace"`
}

const defaultNamespaced = false
const defaultNamespace = "default"

func LoadPluginSettings(source backend.DataSourceInstanceSettings) (*PluginSettings, error) {
	if source.JSONData == nil || len(source.JSONData) < 1 {
		// If no settings have been saved return default values
		return &PluginSettings{
			DefaultNamespaced: defaultNamespaced,
			DefaultNamespace:  defaultNamespace,
		}, nil
	}

	settings := PluginSettings{
		DefaultNamespaced: defaultNamespaced,
		DefaultNamespace:  defaultNamespace,
	}

	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}

	return &settings, nil
}
