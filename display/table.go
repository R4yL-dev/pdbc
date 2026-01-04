package display

import (
	"fmt"
	"pdbc/models"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[96m"  // Platinum - Excellent
	colorYellow = "\033[93m"  // Gold - Très bon
	colorWhite  = "\033[97m"  // Silver - Bon
	colorOrange = "\033[33m"  // Bronze - Moyen
	colorRed    = "\033[91m"  // Borked/Inadequate - Mauvais
	colorGray   = "\033[90m"  // Pending/Unknown - Indéterminé
	colorGreen  = "\033[92m"  // Strong - Confiance forte
	colorDim    = "\033[2;37m" // Weak/Low - Confiance faible
)

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if s == "" || s == "Unknown" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// colorTier returns the tier string with appropriate color
func colorTier(tier string) string {
	tierLower := strings.ToLower(tier)
	capitalized := capitalizeFirst(tier)

	switch tierLower {
	case "platinum":
		return colorCyan + capitalized + colorReset
	case "gold":
		return colorYellow + capitalized + colorReset
	case "silver":
		return colorWhite + capitalized + colorReset
	case "bronze":
		return colorOrange + capitalized + colorReset
	case "borked":
		return colorRed + capitalized + colorReset
	case "pending":
		return colorGray + capitalized + colorReset
	case "unknown":
		return colorGray + capitalized + colorReset
	default:
		return capitalized
	}
}

// colorConfidence returns the confidence string with appropriate color
func colorConfidence(confidence string) string {
	confLower := strings.ToLower(confidence)
	capitalized := capitalizeFirst(confidence)

	switch confLower {
	case "strong":
		return colorGreen + capitalized + colorReset
	case "good", "moderate":
		return colorYellow + capitalized + colorReset
	case "weak", "low":
		return colorDim + capitalized + colorReset
	case "inadequate":
		return colorRed + capitalized + colorReset
	case "unknown":
		return colorGray + capitalized + colorReset
	default:
		return capitalized
	}
}

// isWideChar returns true if the rune takes 2 columns in the terminal
// Detects CJK (Chinese, Japanese, Korean) characters and emojis
func isWideChar(r rune) bool {
	// CJK and other wide character ranges
	return (r >= 0x1100 && r <= 0x115F) || // Hangul Jamo
		(r >= 0x2329 && r <= 0x232A) || // Left/Right-Pointing Angle Brackets
		(r >= 0x2E80 && r <= 0x9FFF) || // CJK Radicals, Hiragana, Katakana, Bopomofo, Hangul, Kanbun, Han
		(r >= 0xAC00 && r <= 0xD7A3) || // Hangul Syllables
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
		(r >= 0xFE10 && r <= 0xFE19) || // Vertical forms
		(r >= 0xFE30 && r <= 0xFE6F) || // CJK Compatibility Forms
		(r >= 0xFF00 && r <= 0xFF60) || // Fullwidth Forms
		(r >= 0xFFE0 && r <= 0xFFE6) || // Fullwidth Forms
		(r >= 0x1F300 && r <= 0x1F9FF) || // Emojis and symbols
		(r >= 0x20000 && r <= 0x2FFFD) || // CJK Extension B-F
		(r >= 0x30000 && r <= 0x3FFFD) // CJK Extension G
}

// displayWidth returns the display width of a string (accounting for wide chars and ANSI codes)
func displayWidth(s string) int {
	// Remove ANSI escape sequences first
	cleaned := s
	for strings.Contains(cleaned, "\033[") {
		start := strings.Index(cleaned, "\033[")
		end := strings.Index(cleaned[start:], "m")
		if end == -1 {
			break
		}
		cleaned = cleaned[:start] + cleaned[start+end+1:]
	}

	// Calculate width accounting for wide characters
	width := 0
	for _, r := range cleaned {
		if isWideChar(r) {
			width += 2
		} else {
			width += 1
		}
	}
	return width
}

// visibleLength is kept for backward compatibility, now uses displayWidth
func visibleLength(s string) int {
	return displayWidth(s)
}

// padRight pads a string to the right with spaces, accounting for ANSI codes
func padRight(s string, width int) string {
	visLen := visibleLength(s)
	if visLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visLen)
}

// PrintTable displays a single table of game results
func PrintTable(results []models.GameResult) {
	if len(results) == 0 {
		fmt.Println("No games found")
		return
	}

	// Calculate column widths (using display width for wide characters)
	maxName := displayWidth("Game Name")
	for _, r := range results {
		titleWidth := displayWidth(r.Title)
		if titleWidth > maxName {
			maxName = titleWidth
		}
	}

	tierWidth := 10
	confidenceWidth := 12

	// Print header
	fmt.Printf("%s | %s | %s\n",
		padRight("Game Name", maxName),
		padRight("Tier", tierWidth),
		padRight("Confidence", confidenceWidth))
	fmt.Printf("%s-+-%s-+-%s\n",
		strings.Repeat("-", maxName),
		strings.Repeat("-", tierWidth),
		strings.Repeat("-", confidenceWidth))

	// Print rows
	for _, r := range results {
		tierColored := colorTier(r.Tier)
		confidenceColored := colorConfidence(r.Confidence)

		fmt.Printf("%s | %s | %s\n",
			padRight(r.Title, maxName),
			padRight(tierColored, tierWidth),
			padRight(confidenceColored, confidenceWidth))
	}
}

// PrintAllResults displays multiple search results in order of searchTerms
func PrintAllResults(searchResults []models.SearchResult, searchTerms []string) {
	// Create map for quick lookup
	resultMap := make(map[string]models.SearchResult)
	for _, sr := range searchResults {
		resultMap[sr.SearchTerm] = sr
	}

	// Print in the original order of search terms
	for i, term := range searchTerms {
		result, exists := resultMap[term]

		// Add blank line before subsequent results (not before first one)
		if i > 0 {
			fmt.Println()
		}

		// Print header for this search term
		fmt.Printf("=== Results for \"%s\" ===\n\n", term)

		if !exists {
			fmt.Println("Error: Failed to retrieve results for this term")
			continue
		}

		if result.Error != nil {
			fmt.Printf("Error: %v\n", result.Error)
			continue
		}

		if len(result.Results) == 0 {
			fmt.Printf("No games found matching '%s'\n", term)
			continue
		}

		// Print the table
		PrintTable(result.Results)
	}
}
