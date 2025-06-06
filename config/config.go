package config

import (
	"go-lucid/node"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	fullNodeConfig *node.FullNodeConfig
)

func MustReadConfig(path string) *node.FullNodeConfig {
	var c *node.FullNodeConfig
	var buf []byte
	var err error

	yamlReader, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer yamlReader.Close()

	buf, err = io.ReadAll(yamlReader)
	if err != nil {
		panic(err)
	}

	c = &node.FullNodeConfig{}
	c.Node.Peers = []string{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		panic(err)
	}

	fullNodeConfig = c

	return c
}

func MustGetFullNodeConfig() *node.FullNodeConfig {
	if fullNodeConfig == nil {
		panic("global config is nil")
	}

	return fullNodeConfig
}
