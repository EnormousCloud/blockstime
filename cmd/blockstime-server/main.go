package main

import (
	"blockstime/internal/config"
	"blockstime/internal/indexer"
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

func main() {
	var cfg config.Config
	readFile(&cfg)
	readEnv(&cfg)

	// fmt.Printf("%+v", cfg)
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	// indexing - save blocks into databases
	for _, network := range cfg.Networks {
		if network.Disabled {
			continue
		}
		ind, err := indexer.New(&network)
		if err != nil {
			log.Fatal(err)
		}
		if err := ind.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *config.Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *config.Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
