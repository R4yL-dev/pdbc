package models

// Game represents a Steam game from the ProtonDB games list
type Game struct {
	AppID int    `json:"appId"`
	Title string `json:"title"`
}

// TierInfo represents ProtonDB tier data for a game
type TierInfo struct {
	Tier             string  `json:"tier"`
	BestReportedTier string  `json:"bestReportedTier"`
	Confidence       string  `json:"confidence"`
	Score            float64 `json:"score"`
	Total            int     `json:"total"`
	TrendingTier     string  `json:"trendingTier"`
}

// GameResult combines game info with tier data for display
type GameResult struct {
	Game
	Tier       string
	Confidence string
}

// SearchResult groups results by search term
type SearchResult struct {
	SearchTerm string
	Results    []GameResult
	Error      error // Pour capturer erreurs dans goroutine
}
