package config

import (
	"github.com/joho/godotenv"
	"os"
	"vaults-operator/utils"
)

var AppConfig *Config

type Config struct {
	GraphQlEndpoint string
}

func init() {
	if os.Getenv("RUNNING_IN_CONTAINER") != "true" {
		if err := godotenv.Load(); err != nil {
			utils.Log.Error("Error loading .env file")
		} else {
			utils.Log.Info(".env file successfully loaded")
		}
	}

	var err error
	AppConfig, err = LoadConfig()
	if err != nil {
		utils.Log.Fatalf("Error loading config: %s", err)
	}

	utils.Log.Debugf("GraphQl endpoint: %s", AppConfig.GraphQlEndpoint)
}

func LoadConfig() (*Config, error) {
	return &Config{
		GraphQlEndpoint: GetStringEnv("GRAPHQL_ENDPOINT", "https://dev-graph.fathom.fi/subgraphs/name/vaults-subgraph"),
	}, nil
}
