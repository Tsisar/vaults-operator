package config

import (
	"os"
	"vaults-operator/utils"
)

func GetStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		utils.Log.Warnf("%s not found in environment variables, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}
