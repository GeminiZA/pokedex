package pokeAPI

import (
	"encoding/json"
	"io"
	"net/http"

	pokecache "github.com/GeminiZA/pokedex/internal/pokeCache"
)

const locationAPIEndpoint = "https://pokeapi.co/api/v2/location-area/"

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
