package server

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	Path     string              `yaml:"path"`
	Method   []string            `yaml:"method"`
	Response string              `yaml:"response"`
	Headers  map[string][]string `yaml:"headers"`
}

type Model struct {
	Port      int        `yaml:"port"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

const CacheMaxItemSizeDefault = int64(1024 * 1024)

func ReadModel(filename string) (*Model, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	model := Model{}
	err = yaml.Unmarshal(data, &model)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	return &model, nil
}
