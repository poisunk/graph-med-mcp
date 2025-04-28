package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Port int `yaml:"port"`

	Neo4j struct {
		Addr     string `yaml:"addr"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"neo4j"`
}

func MustLoad(path string, v any) {
	if err := Load(path, v); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}

func Load(file string, v any) error {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &v)
	return err
}
