package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"vaults-operator/graph"
)

type Strategy struct {
	ID         string `json:"id"`
	AssetsPart string `json:"assetsPart,omitempty"`
}

type Vault struct {
	ID            string     `json:"id"`
	UseAssetsPart string     `json:"useAssetsPart,omitempty"`
	Strategies    []Strategy `json:"strategies"`
}

type Data struct {
	Vaults []Vault `json:"vaults"`
}

type Addresses struct {
	Data Data `json:"data"`
}

func readAndModifyJSON(filePath string) (*Addresses, error) {
	bytes, err := graph.ExecuteQuery()
	if err != nil {
		return nil, err
	}

	var addresses Addresses
	if err := json.Unmarshal(bytes, &addresses); err != nil {
		return nil, err
	}

	// Modify the data

	return &addresses, nil
}

func main() {
	// Create a new Gin server
	r := gin.Default()
	r.GET("/data-json", func(c *gin.Context) {
		modifiedData, err := readAndModifyJSON("data.json") // Замініть на свій файл
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, modifiedData)
	})

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
