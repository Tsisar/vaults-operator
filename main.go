package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
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

// Функція для читання та модифікації JSON
func readAndModifyJSON(filePath string) (*Addresses, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var addresses Addresses
	if err := json.Unmarshal(bytes, &addresses); err != nil {
		return nil, err
	}

	// Тут можна додати логіку модифікації JSON
	// Наприклад, додати assetsPart до кожної стратегії

	return &addresses, nil
}

func main() {
	// Створення серверу і роута
	r := gin.Default()
	r.GET("/data-json", func(c *gin.Context) {
		modifiedData, err := readAndModifyJSON("data.json") // Замініть на свій файл
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, modifiedData)
	})

	// Запуск серверу
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}