package main

import (
	"fmt"
	"os"
	"pdbc/api"
	"pdbc/display"
	"pdbc/models"
	"pdbc/spinner"
	"sync"
)

func main() {
	// Parse CLI arguments
	searchTerms := os.Args[1:]

	if len(searchTerms) == 0 {
		printUsage()
		os.Exit(1)
	}

	// Create spinner for loading indication
	spin := spinner.New("Loading datas")
	spin.Start()

	// Create channel for results and WaitGroup for synchronization
	resultsChan := make(chan models.SearchResult, len(searchTerms))
	var wg sync.WaitGroup

	// Launch a goroutine for each search term
	for _, term := range searchTerms {
		wg.Add(1)
		go searchGames(term, resultsChan, &wg)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	results := make([]models.SearchResult, 0, len(searchTerms))
	for result := range resultsChan {
		results = append(results, result)
	}

	// Stop spinner and clear the line
	spin.Stop()

	// Display results in original order
	display.PrintAllResults(results, searchTerms)
}

// searchGames performs a search for a single term and sends results to channel
func searchGames(term string, resultsChan chan<- models.SearchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// Search games on Steam Store
	matchedGames, err := api.SearchGames(term)
	if err != nil {
		// Send error result
		resultsChan <- models.SearchResult{
			SearchTerm: term,
			Results:    nil,
			Error:      err,
		}
		return
	}

	// Fetch tier data for matched games concurrently
	gameResults := fetchTiersForGames(matchedGames)

	// Send results to channel
	resultsChan <- models.SearchResult{
		SearchTerm: term,
		Results:    gameResults,
		Error:      nil,
	}
}

// fetchTiersForGames fetches tier information for multiple games concurrently
// Uses a semaphore to limit concurrent HTTP requests
func fetchTiersForGames(games []models.Game) []models.GameResult {
	if len(games) == 0 {
		return []models.GameResult{}
	}

	results := make([]models.GameResult, len(games))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Max 10 concurrent HTTP requests

	for i, game := range games {
		wg.Add(1)
		go func(idx int, g models.Game) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }() // Release semaphore

			tier, err := api.FetchGameTier(g.AppID)
			if err != nil {
				// Mark as Unknown on error
				results[idx] = models.GameResult{
					Game:       g,
					Tier:       "Unknown",
					Confidence: "Unknown",
				}
			} else {
				results[idx] = models.GameResult{
					Game:       g,
					Tier:       tier.Tier,
					Confidence: tier.Confidence,
				}
			}
		}(i, game)
	}

	wg.Wait()
	return results
}

func printUsage() {
	fmt.Println("Usage: pdbc <search_term> [search_term2] [search_term3] ...")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  pdbc \"Anno\"")
	fmt.Println("  pdbc \"Half-Life\"")
	fmt.Println("  pdbc \"Anno\" \"Cyberpunk\" \"Half-Life\"")
}
