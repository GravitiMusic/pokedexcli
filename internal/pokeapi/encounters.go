package pokeapi

import (
	"encoding/json"
	"net/http"
)

type Pokemon struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
}

type EncountersResponse struct {
	Encounters []Pokemon `json:"pokemon_encounters"`
}

func (c *Client) GetEncounters(url string) ([]string, error) {
	if cached, ok := c.cache.Get(url); ok {
		var result EncountersResponse
		if err := json.Unmarshal(cached, &result); err != nil {
			return nil, err
		}
		var names []string
		for _, encounter := range result.Encounters {
			names = append(names, encounter.Pokemon.Name)
		}
		return names, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result EncountersResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var names []string
	for _, encounter := range result.Encounters {
		names = append(names, encounter.Pokemon.Name)
	}

	return names, nil
}