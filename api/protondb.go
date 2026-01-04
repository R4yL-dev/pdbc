package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"pdbc/models"
	"time"
)

const (
	steamSearchURL = "https://store.steampowered.com/api/storesearch/?term=%s&l=english&cc=US"
	tierURLBase    = "https://www.protondb.com/api/v1/reports/summaries/%d.json"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// SteamSearchResponse represents the response from Steam Store search API
type SteamSearchResponse struct {
	Items []SteamSearchItem `json:"items"`
	Total int               `json:"total"`
}

// SteamSearchItem represents a single game in Steam Store search results
type SteamSearchItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// SearchGames searches for games on Steam Store by name
func SearchGames(searchTerm string) ([]models.Game, error) {
	// URL encode the search term
	encodedTerm := url.QueryEscape(searchTerm)
	searchURL := fmt.Sprintf(steamSearchURL, encodedTerm)

	resp, err := httpClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search Steam for '%s': %w", searchTerm, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Steam search API returned status %d", resp.StatusCode)
	}

	var searchResp SteamSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode Steam search JSON: %w", err)
	}

	// Convert Steam search items to our Game model
	games := make([]models.Game, 0, len(searchResp.Items))
	for _, item := range searchResp.Items {
		// Only include apps (not DLC, videos, etc.)
		if item.Type == "app" {
			games = append(games, models.Game{
				AppID: item.ID,
				Title: item.Name,
			})
		}
	}

	return games, nil
}

// FetchGameTier retrieves tier information for a specific game by AppID
func FetchGameTier(appID int) (*models.TierInfo, error) {
	url := fmt.Sprintf(tierURLBase, appID)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tier for appID %d: %w", appID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tier API returned status %d for appID %d", resp.StatusCode, appID)
	}

	var tierInfo models.TierInfo
	if err := json.NewDecoder(resp.Body).Decode(&tierInfo); err != nil {
		return nil, fmt.Errorf("failed to decode tier JSON for appID %d: %w", appID, err)
	}

	return &tierInfo, nil
}
