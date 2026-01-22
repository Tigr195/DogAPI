package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DogResponse struct {
	Message []string `json:"message"`
	Status  string   `json:"status"`
}

type DogClient struct {
	apiKey string
	client *http.Client
}

type DogItem struct {
	Id  string `json:"id"`
	URL string `json:"url"`
}

func NewDogClient() *DogClient {
	return &DogClient{
		apiKey: os.Getenv("API_KEY"),
		client: &http.Client{},
	}
}

func (c *DogClient) GetDogs(limit int) ([]DogItem, error) {
	req, err := http.NewRequest("GET", "https://api.thedogapi.com/v1/images/search", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("limit", fmt.Sprintf("%d", limit))

	req.URL.RawQuery = q.Encode()

	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items []DogItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}
