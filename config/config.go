package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	WikiPath   string `yaml:"wikiPath"`
	LlmUrl     string `yaml:"llmUrl"`
	LlmModel   string `yaml:"llmModel"`
	Prometheus struct {
		Tg struct {
			URL      string `yaml:"url"`
			Name     string `yaml:"name"`
			Password string `yaml:"password"`
		} `yaml:"tg"`
	} `yaml:"prometheus"`
}

func Init() *Conf {
	c := &Conf{}
	bytes, err := os.ReadFile("config/app.yaml")
	if err != nil {
		fmt.Println("read app.yaml err:", err)
		panic(err)
	}
	if err = yaml.Unmarshal(bytes, c); err != nil {
		fmt.Println("parse app.yaml err:", err)
		panic(err)
	}
	return c
}
