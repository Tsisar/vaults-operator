package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vaults-operator/config"
	"vaults-operator/utils"
)

var query = `{"query": "{ vaults { id, strategies { id } } }"}`

type Strategy struct {
	ID         string `json:"id"`
	AssetsPart string `json:"assetsPart,omitempty"`
}

type Vault struct {
	ID         string     `json:"id"`
	Strategies []Strategy `json:"strategies"`
}

type Data struct {
	Vaults []Vault `json:"vaults"`
}

type Addresses struct {
	Data Data `json:"data"`
}

func ParseAddresses(data []byte) (*Addresses, error) {
	var response Addresses
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &response, nil
}

func ExecuteQuery() ([]byte, error) {
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
