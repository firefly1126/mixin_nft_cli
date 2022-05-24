package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseUrl string `yaml:"base_url"`
	Token   string `yaml:"token"`
	UserID  string `yaml:"user_id"`
}

func Load(path string, c *Config) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		panic(err)
	}
}
