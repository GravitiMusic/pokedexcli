package pokeapi

import (
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
)

const PokemonURL = "https://pokeapi.co/api/v2/pokemon/"

type Pokedex struct {
	Caught map[string]Pokemon
}

type Pokemon struct {
	Id			   int    `json:"id"`
	Name		   string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height		   int    `json:"height"`
	Weight		   int    `json:"weight"`
	Types		   []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Stats 		   []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
}

type Encounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
}

type EncountersResponse struct {
	Encounters []Encounter `json:"pokemon_encounters"`
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

	if data, err := json.Marshal(result); err == nil {
		c.cache.Add(url, data)
	}

	return names, nil
}

func (c *Client) CatchPokemon(url, name string, dex *Pokedex) (int, error) {
	if cached, ok := c.cache.Get(url); ok {
		var result Pokemon
		if err := json.Unmarshal(cached, &result); err != nil {
			return 0, err
		}
		probability := (0.05 + 0.85 * math.Pow((350.0 - float64(result.BaseExperience)) / 315.0, 1.5)) * 100
		if rand.Intn(100) + 1 <= int(probability) {
			dex.Caught[name] = result
			return 1, nil
		}
		return 0, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	
	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var result Pokemon
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, err
	}

	if data, err := json.Marshal(result); err == nil {
		c.cache.Add(url, data)
	}

	probability := (0.05 + 0.85 * math.Pow((350.0 - float64(result.BaseExperience)) / 315.0, 1.5)) * 100

	if rand.Intn(100) + 1 <= int(probability) {
		dex.Caught[name] = result
		return 1, nil
	}

	return 0, nil
}