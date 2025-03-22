package config

import (
	"go-lucid/node"
	"io"
	"os"

	"gopkg.in/yaml.v3"
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

	return c
}
