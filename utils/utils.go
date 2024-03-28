package utils

import (
	"os"
)

func Unhealthy() {
	err := os.Remove("/tmp/healthy")
	if err != nil {
		Log.Errorf("Error removing /tmp/healthy file: %s", err)
	} else {
		Log.Info("The /tmp/healthy file was successfully removed.")
	}
}
