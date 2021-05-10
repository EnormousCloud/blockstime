package config

import (
	"blockstime/internal/engines"
	"errors"
	"fmt"
	"strings"
)

type Network struct {
	// Name of the network
	Name string `json:"name" yaml:"name"`
	// Engine of the network. Ethereum ENM is the only supported so far
	Engine string `json:"engine" yaml:"engine"`
	// Nodes are RPC-hosts to be used for reading networks
	Nodes []engines.NodeRPC `json:"nodes" yaml:"nodes"`
	// In case of local file storage - name of the local file
	LocalPath string `json:"localpath" yaml:"localpath"`
	// Whether this network is disabled for indexing
	Disabled bool `json:"disabled" yaml:"disabled"`
}

func (n *Network) Validate() error {
	if len(n.Name) == 0 {
		return errors.New("name of network is not provided")
	}
	if len(n.Engine) > 0 && strings.ToLower(n.Engine) != "evm" {
		return errors.New("invalid network engine. Only EVM is allowed so far")
	}
	if len(n.Nodes) == 0 {
		return errors.New("no nodes in the network")
	}
	if len(n.LocalPath) == 0 {
		return errors.New("storage localpath is not specified")
	}
	return nil
}

type Config struct {
	// List of networks to be parsed
	Networks []Network `json:"networks" yaml:"networks"`
}

func (c *Config) Validate() error {
	if len(c.Networks) == 0 {
		return errors.New("no networks in configuration")
	}
	for index, n := range c.Networks {
		if err := n.Validate(); err != nil {
			return fmt.Errorf("network %d: %v", index+1, err.Error())
		}
	}
	return nil
}
