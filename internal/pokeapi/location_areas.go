package pokeapi

import (
	"encoding/json"
	"net/http"
	"time"

	"pokedexcli/internal/pokecache"
)

const LocationAreaURL = "https://pokeapi.co/api/v2/location-area/"

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type Client struct {
	cache      *pokecache.Cache
	httpClient *http.Client
}

func NewClient(cacheInterval time.Duration) *Client {
	return &Client{
		cache:      pokecache.NewCache(cacheInterval),
		httpClient: &http.Client{},
	}
}

func (c *Client) GetLocationAreas(url string) (LocationAreaResponse, error) {
	if cached, ok := c.cache.Get(url); ok {
		var result LocationAreaResponse
		if err := json.Unmarshal(cached, &result); err != nil {
			return LocationAreaResponse{}, err
		}
		return result, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer res.Body.Close()

	var result LocationAreaResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return LocationAreaResponse{}, err
	}

	if data, err := json.Marshal(result); err == nil {
		c.cache.Add(url, data)
	}

	return result, nil
}
