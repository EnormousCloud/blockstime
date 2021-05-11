package main

import (
	"blockstime/internal/config"
	"blockstime/internal/indexer"
	"blockstime/internal/server"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	cli "github.com/jawher/mow.cli"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var app = cli.App("blockstime-server", "Blockstime API")

var (
	cfgFile = app.String(cli.StringOpt{
		Name:   "config",
		Desc:   "config.yml full file path",
		EnvVar: "CONFIG_FILE",
		Value:  "./config.yml",
	})
	networkToIndex = app.String(cli.StringOpt{
		Name:  "index",
		Desc:  "network from config.yml to be indexed (HTTP server will not start)",
		Value: "",
	})
	netListen = app.String(cli.StringOpt{
		Name:  "listen",
		Desc:  "TCP address to be listened",
		Value: "0.0.0.0:8080",
	})
)

// @title Blockstime API
// @version 1.0
// @description Microservice to seamlessly convert time of blocks into Unix timestamps
// @host localhost:8080
// @BasePath /api
func main() {
	app.Action = startApp
	app.Run(os.Args)
}

func startApp() {
	var cfg config.Config
	readFile(&cfg, *cfgFile)
	readEnv(&cfg)

	// fmt.Printf("%+v", cfg)
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	if len(*networkToIndex) > 0 {
		log.Printf("[main] Indexing %v", *networkToIndex)
		// indexing - save blocks into databases
		for _, network := range cfg.Networks {
			// if network.Disabled {
			// 	continue
			// }
			if network.Name == *networkToIndex {
				ind, err := indexer.New(&network)
				if err != nil {
					log.Fatal(err)
				}
				if err := ind.Run(); err != nil {
					log.Fatal(err)
				}
			}
		}
		log.Fatalf("Failed to index %v - network not found in configuration",
			*networkToIndex)
	}
	// start a web server otherwise
	r := gin.Default()
	c := server.NewController()
	apiV1 := r.Group("/api")
	{
		apiV1.GET("ping", c.Ping)
		apiV1.GET("blocks", c.GetBlocksFromPeriods)
		apiV1.GET("periods", c.GetPeriodFromBlocks)
		apiV1.GET("stats/daily", c.GetStatsDaily)
		apiV1.GET("stats/yearly", c.GetStatsYearly)
	}
	log.Printf("Starting server %v\n", *netListen)
	r.Run(*netListen) // listen and serve on
}

func processError(err error) {
	log.Println(err)
	os.Exit(2)
}

func readFile(cfg *config.Config, cfgFile string) {
	f, err := os.Open(cfgFile)
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
