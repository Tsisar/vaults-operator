package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"vaults-operator/config"
	"vaults-operator/utils"
)

var query = `{"query": "{ vaults { id, strategies { id } } }"}`

type Strategy struct {
	ID         string  `json:"id"`
	AssetsPart float64 `json:"assetsPart,omitempty"`
}

type Vault struct {
	ID            string     `json:"id"`
	UseAssetsPart bool       `json:"useAssetsPart,omitempty"`
	Strategies    []Strategy `json:"strategies"`
}

type Data struct {
	Vaults []Vault `json:"vaults"`
}

type Addresses struct {
	Data Data `json:"data"`
}

// Global variable to store the data, for demonstration
var db Addresses

func readJson() (*Addresses, error) {
	data, err := executeQuery()
	//data, err := ioutil.ReadFile("data.json") // Test data
	if err != nil {
		return nil, err
	}

	var addresses Addresses
	if err = json.Unmarshal(data, &addresses); err != nil {
		return nil, err
	}

	for i, vault := range addresses.Data.Vaults {
		for _, dbVault := range db.Data.Vaults {
			if vault.ID == dbVault.ID {
				addresses.Data.Vaults[i].UseAssetsPart = dbVault.UseAssetsPart
				for j, strategy := range vault.Strategies {
					for _, dbStrategy := range dbVault.Strategies {
						if strategy.ID == dbStrategy.ID {
							addresses.Data.Vaults[i].Strategies[j].AssetsPart = dbStrategy.AssetsPart
							break
						}
					}
				}
				break
			}
		}
	}

	return &addresses, nil
}

func main() {
	// Create a new Gin server
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/data-json", func(c *gin.Context) {
		data, err := readJson()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	r.POST("/update-json", func(c *gin.Context) {
		var updatedData Addresses

		if err := c.BindJSON(&updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		updateVaultsData(updatedData)

		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data updated successfully"})
	})

	r.GET("/", func(c *gin.Context) {
		data, err := readJson()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Vaults": data.Data.Vaults,
		})
	})

	// Run the server
	port := fmt.Sprintf(":%s", config.AppConfig.Port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}

func updateVaultsData(newData Addresses) {
	db = newData
}

func executeQuery() ([]byte, error) {
	resp, err := http.Post(config.AppConfig.GraphQlEndpoint, "application/json", bytes.NewBufferString(query))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Log.Errorf("Error closing response body: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
