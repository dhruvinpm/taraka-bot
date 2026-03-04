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

type Hunter struct {
	store *memory.Store
	pool  *browser.BrowserPool
}

func NewHunter(store *memory.Store, pool *browser.BrowserPool) *Hunter {
	return &Hunter{store: store, pool: pool}
}

func (h *Hunter) Hunt(ctx context.Context, cfg HuntConfig) ([]*memory.Lead, error) {
	var all []*memory.Lead

	query := fmt.Sprintf("%s importer buyer %s", cfg.Niche, cfg.Country)

	googleHunter := NewGoogleHunter()
	gLeads, err := googleHunter.SearchGoogle(query, cfg.Country)
	if err != nil {
		log.Printf("google hunt error: %v", err)
	} else {
		all = append(all, gLeads...)
	}

	ddgHunter := NewDuckDuckGoHunter()
	dLeads, err := ddgHunter.SearchDDG(query)
	if err != nil {
		log.Printf("ddg hunt error: %v", err)
	} else {
		all = append(all, dLeads...)
	}

	dirHunter := NewDirectoryHunter()
	dirLeads, err := dirHunter.SearchDirectories(query, cfg.Country)
	if err != nil {
		log.Printf("directory hunt error: %v", err)
	} else {
		all = append(all, dirLeads...)
	}

	if cfg.Mode == Moderate || cfg.Mode == Aggressive {
		redditHunter := NewRedditHunter()
		rLeads, err := redditHunter.SearchReddit(cfg.Niche)
		if err != nil {
			log.Printf("reddit hunt error: %v", err)
		} else {
			all = append(all, rLeads...)
		}
	}

	if cfg.MaxLeads > 0 && len(all) > cfg.MaxLeads {
		all = all[:cfg.MaxLeads]
	}

	for _, lead := range all {
		if err := h.store.UpsertLead(lead); err != nil {
			log.Printf("upsert lead error: %v", err)
		}
	}

	return all, nil
}
