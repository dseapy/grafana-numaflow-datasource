package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type PluginSettings struct {
	Namespaced bool   `json:"namespaced"`
	Namespace  string `json:"namespace"`
}

func LoadPluginSettings(source backend.DataSourceInstanceSettings) (*PluginSettings, error) {
	settings := PluginSettings{
		Namespaced: false,
		Namespace:  "default",
	}

	if source.JSONData == nil || len(source.JSONData) < 1 {
		backend.Logger.Debug("No settings found, using default settings")
		return &settings, nil
	}

	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}
	backend.Logger.Debug("Successfully parsed settings", "namespaced", settings.Namespaced, "namespace", settings.Namespace)

	return &settings, nil
}
