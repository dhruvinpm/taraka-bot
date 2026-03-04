package hunter

import (
	"context"
	"fmt"
	"log"

	"github.com/dhruvinpm/taraka-bot/browser"
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

// HuntResult holds per-source lead counts and the aggregated lead list.
type HuntResult struct {
	Leads          []*memory.Lead
	GoogleCount    int
	DDGCount       int
	LinkedInCount  int
	DirectoryCount int
	RedditCount    int
	TotalLeads     int
	EmailCount     int
	Errors         []string
}

type Hunter struct {
	store *memory.Store
	pool  *browser.BrowserPool
}

func NewHunter(store *memory.Store, pool *browser.BrowserPool) *Hunter {
	return &Hunter{store: store, pool: pool}
}

func (h *Hunter) Hunt(ctx context.Context, cfg HuntConfig) (*HuntResult, error) {
	result := &HuntResult{}

	query := fmt.Sprintf("%s importer buyer %s", cfg.Niche, cfg.Country)

	googleHunter := NewGoogleHunter(h.pool)
	gLeads, err := googleHunter.SearchGoogle(query, cfg.Country)
	if err != nil {
		log.Printf("google hunt error: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("google: %v", err))
	} else {
		result.GoogleCount = len(gLeads)
		result.Leads = append(result.Leads, gLeads...)
	}

	ddgHunter := NewDuckDuckGoHunter()
	dLeads, err := ddgHunter.SearchDDG(query)
	if err != nil {
		log.Printf("ddg hunt error: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("duckduckgo: %v", err))
	} else {
		result.DDGCount = len(dLeads)
		result.Leads = append(result.Leads, dLeads...)
	}

	linkedInHunter := NewLinkedInHunter(h.pool)
	lLeads, err := linkedInHunter.SearchLinkedIn(fmt.Sprintf("FMCG distributor %s", cfg.Country))
	if err != nil {
		log.Printf("linkedin hunt error: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("linkedin: %v", err))
	} else {
		result.LinkedInCount = len(lLeads)
		result.Leads = append(result.Leads, lLeads...)
	}

	dirHunter := NewDirectoryHunter()
	dirLeads, err := dirHunter.SearchDirectories(query, cfg.Country)
	if err != nil {
		log.Printf("directory hunt error: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("directories: %v", err))
	} else {
		result.DirectoryCount = len(dirLeads)
		result.Leads = append(result.Leads, dirLeads...)
	}

	if cfg.Mode == Moderate || cfg.Mode == Aggressive {
		redditHunter := NewRedditHunter()
		rLeads, err := redditHunter.SearchReddit(cfg.Niche)
		if err != nil {
			log.Printf("reddit hunt error: %v", err)
			result.Errors = append(result.Errors, fmt.Sprintf("reddit: %v", err))
		} else {
			result.RedditCount = len(rLeads)
			result.Leads = append(result.Leads, rLeads...)
		}
	}

	if cfg.MaxLeads > 0 && len(result.Leads) > cfg.MaxLeads {
		result.Leads = result.Leads[:cfg.MaxLeads]
	}

	for _, lead := range result.Leads {
		if lead.Email != "" {
			result.EmailCount++
		}
		if err := h.store.UpsertLead(lead); err != nil {
			log.Printf("upsert lead error: %v", err)
		}
	}

	result.TotalLeads = len(result.Leads)
	return result, nil
}
