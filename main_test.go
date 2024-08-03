package main

import (
	"gopkg.in/yaml.v2"
	"os"
	"sort"
	"testing"
)

func TestMarshalConfig(t *testing.T) {
	cfg := readConfig(configPath)
	sort.Slice(cfg.SearchList, func(i, j int) bool {
		return cfg.SearchList[i] < cfg.SearchList[j]
	})

	yamlBytes, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(yamlBytes)
	if err != nil {
		panic(err)
	}
}
