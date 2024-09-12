package pokeAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	pokecache "github.com/GeminiZA/pokedex/internal/pokeCache"
)

const (
	locationAPIEndpoint = "https://pokeapi.co/api/v2/location-area/"
)

type MapRes struct {
	Count   int     `json:"count"`
	Next    *string `json:"next"`
	Prev    *string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func parseMapRes(data []byte) (*MapRes, error) {
	var mapRes MapRes
	err := json.Unmarshal(data, &mapRes)
	if err != nil {
		return nil, err
	}
	return &mapRes, nil
}

type API struct {
	ApiCache *pokecache.Cache
}

func (api *API) Map(url string) (*MapRes, error) {
	endpoint := url
	if endpoint == "" {
		endpoint = locationAPIEndpoint
	}
	if cachedRes, ok := api.ApiCache.Get(endpoint); ok {
		return parseMapRes(cachedRes)
	}
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	api.ApiCache.Add(endpoint, data)
	mapRes, err := parseMapRes(data)
	if err != nil {
		return nil, err
	}
	return mapRes, nil
}

type ExploreRes struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func parseExloreRes(data []byte) (*ExploreRes, error) {
	var exploreRes ExploreRes
	err := json.Unmarshal(data, &exploreRes)
	if err != nil {
		return nil, err
	}
	return &exploreRes, nil
}

func (api *API) Explore(location string) (*ExploreRes, error) {
	if location == "" {
		return nil, fmt.Errorf("no location to be passed")
	}
	endpoint, err := url.JoinPath(locationAPIEndpoint, location)
	if err != nil {
		return nil, err
	}
	if cachedRes, ok := api.ApiCache.Get(endpoint); ok {
		return parseExloreRes(cachedRes)
	}
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, nil
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if string(data) == "Not Found" {
		return nil, fmt.Errorf("area not found")
	}
	api.ApiCache.Add(endpoint, data)
	exploreRes, err := parseExloreRes(data)
	if err != nil {
		return nil, err
	}
	return exploreRes, nil
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
}

const pokemonAPIEndpoint = "https://pokeapi.co/api/v2/pokemon/"

func parsePokemon(data []byte) (*Pokemon, error) {
	var pokemon Pokemon
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return nil, err
	}
	return &pokemon, nil
}

func (api *API) GetPokemon(name string) (*Pokemon, error) {
	if name == "" {
		return nil, fmt.Errorf("no name passed")
	}
	endpoint, err := url.JoinPath(pokemonAPIEndpoint, name)
	if err != nil {
		return nil, err
	}
	if cachedRes, ok := api.ApiCache.Get(endpoint); ok {
		return parsePokemon(cachedRes)
	}
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, nil
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if string(data) == "Not Found" {
		return nil, fmt.Errorf("pokemon not found")
	}
	api.ApiCache.Add(endpoint, data)
	pokemon, err := parsePokemon(data)
	if err != nil {
		return nil, err
	}
	return pokemon, nil
}
