package config

import (
	"encoding/json"
	"fmt"
	"log"
)

type SwaggerConfig struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Host        string `json:"host"`
	BasePath    string `json:"BasePath"`
	Schema      string `json:"schema"`
}

func newSwaggerConfig(cfg SwaggerConfig) *SwaggerConfig {
	return &SwaggerConfig{
		Title:       cfg.Title,
		Description: cfg.Description,
		Version:     cfg.Version,
		Host:        cfg.Host,
		BasePath:    cfg.BasePath,
		Schema:      cfg.Schema,
	}
}

func ReadFromFileSwagger(content []byte) (SwaggerConfig, error) {
	var swaggerCfgSchema SwaggerConfig

	if err := json.Unmarshal(content, &swaggerCfgSchema); err != nil {
		return SwaggerConfig{}, fmt.Errorf("cannot parse swagger-config: %w", err)
	}

	jsonLog, err := json.Marshal(struct {
		JSON SwaggerConfig `json:"json"`
	}{swaggerCfgSchema})
	if err != nil {
		log.Printf("failed to marshal config to JSON: %+v", err)
	}

	log.Printf("%s", jsonLog)

	newSwaggerConfig := newSwaggerConfig(swaggerCfgSchema)

	return *newSwaggerConfig, nil
}
