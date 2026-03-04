package hunter

import (
	"context"
	"fmt"
	"log"

	"github.com/dhruvinpm/taraka-bot/browser"
	"github.com/dhruvinpm/taraka-bot/llm"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type AggressionMode string

const (
	Gentle     AggressionMode = "gentle"
	Moderate   AggressionMode = "moderate"
	Aggressive AggressionMode = "aggressive"
)

type HuntConfig struct {
	Country  string
	Niche    string
	Mode     AggressionMode
	MaxLeads int
}

// HuntResult carries both the found leads and any per-source errors so callers
// can report them to the user (e.g. via Telegram) rather than silently logging.
type HuntResult struct {
	Leads  []*memory.Lead
	Errors []string
}

type Hunter struct {
	store    *memory.Store
	pool     *browser.BrowserPool
	llmProv  llm.Provider
}

func NewHunter(store *memory.Store, pool *browser.BrowserPool, llmProvider llm.Provider) *Hunter {
	return &Hunter{store: store, pool: pool, llmProv: llmProvider}
}

func (h *Hunter) Hunt(ctx context.Context, cfg HuntConfig) (*HuntResult, error) {
	result := &HuntResult{}

	query := fmt.Sprintf("%s importer buyer %s", cfg.Niche, cfg.Country)

	// Primary source: Gemini AI (most reliable — not blocked by anti-bot measures).
	if h.llmProv != nil {
		geminiHunter := NewGeminiHunter(h.llmProv)
		gLeads, err := geminiHunter.SearchGemini(ctx, query, cfg.Country, cfg.MaxLeads)
		if err != nil {
			log.Printf("gemini hunt error: %v", err)
			result.Errors = append(result.Errors, fmt.Sprintf("Gemini: %v", err))
		} else {
			result.Leads = append(result.Leads, gLeads...)
		}
	} else {
		result.Errors = append(result.Errors, "Gemini: not configured (no API key)")
	}

	// Fallback scraper sources — these often fail due to anti-bot protections,
	// but errors are now collected and surfaced to the user.
	googleHunter := NewGoogleHunter()
	gLeads, err := googleHunter.SearchGoogle(query, cfg.Country)
	if err != nil {
		log.Printf("google hunt error: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("Google: %v", err))
	} else {
		result.Leads = append(result.Leads, gLeads...)
	}

	ddgHunter := NewDuckDuckGoHunter()
	dLeads, err := ddgHunter.SearchDDG(query)
	if err != nil {
		log.Printf("ddg hunt error: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("DuckDuckGo: %v", err))
	} else {
		result.Leads = append(result.Leads, dLeads...)
	}

	dirHunter := NewDirectoryHunter()
	dirLeads, err := dirHunter.SearchDirectories(query, cfg.Country)
	if err != nil {
		log.Printf("directory hunt error: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("Directories: %v", err))
	} else {
		result.Leads = append(result.Leads, dirLeads...)
	}

	if cfg.Mode == Moderate || cfg.Mode == Aggressive {
		redditHunter := NewRedditHunter()
		rLeads, err := redditHunter.SearchReddit(cfg.Niche)
		if err != nil {
			log.Printf("reddit hunt error: %v", err)
			result.Errors = append(result.Errors, fmt.Sprintf("Reddit: %v", err))
		} else {
			result.Leads = append(result.Leads, rLeads...)
		}
	}

	if cfg.MaxLeads > 0 && len(result.Leads) > cfg.MaxLeads {
		result.Leads = result.Leads[:cfg.MaxLeads]
	}

	for _, lead := range result.Leads {
		if err := h.store.UpsertLead(lead); err != nil {
			log.Printf("upsert lead error: %v", err)
		}
	}

	return result, nil
}
