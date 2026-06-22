package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type PlaceSuggestion struct {
	Title     string  `json:"title"`
	Address   string  `json:"address"`
	Province  string  `json:"province"`
	City      string  `json:"city"`
	District  string  `json:"district"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey, httpClient: &http.Client{}}
}

func (c *Client) PlaceSearch(ctx context.Context, keyword string) ([]PlaceSuggestion, error) {
	if c.apiKey == "" {
		return []PlaceSuggestion{}, nil
	}

	baseURL := "https://apis.map.qq.com/ws/place/v1/suggestion"
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("key", c.apiKey)
	params.Set("region", "全国")

	resp, err := c.httpClient.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("place search: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Status int `json:"status"`
		Data   []struct {
			Title    string `json:"title"`
			Address  string `json:"address"`
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if result.Status != 0 {
		return nil, fmt.Errorf("place search failed: status=%d", result.Status)
	}

	suggestions := make([]PlaceSuggestion, 0, len(result.Data))
	for _, d := range result.Data {
		suggestions = append(suggestions, PlaceSuggestion{
			Title:     d.Title,
			Address:   d.Address,
			Province:  d.Province,
			City:      d.City,
			District:  d.District,
			Latitude:  d.Location.Lat,
			Longitude: d.Location.Lng,
		})
	}
	return suggestions, nil
}

func GenerateNavigationURL(lat, lng float64, address string) string {
	return fmt.Sprintf("https://uri.amap.com/navigation?to=%f,%f,%s", lng, lat, url.QueryEscape(address))
}
