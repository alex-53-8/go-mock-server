package server

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Mapping struct {
	Path     string   `yaml:"path"`
	Method   []string `yaml:"method"`
	Response string   `yaml:"response"`
	Headers  []Header `yaml:"headers"`
}

type Cache struct {
	MaxItemSize int64 `yaml:"maxItemSize"`
}

type Model struct {
	Port     int       `yaml:"port"`
	Mappings []Mapping `yaml:"mappings"`
	Cache    Cache     `yaml:"cache"`
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

	if model.Cache.MaxItemSize == 0 {
		model.Cache.MaxItemSize = CacheMaxItemSizeDefault
	}

	return &model, nil
}
