package analyst

import (
	"log"

	"github.com/dhruvinpm/taraka-bot/memory"
)

type Analyst struct {
	Store    *memory.Store
	verifier *EmailVerifier
	scanner  *WebsiteScanner
	matcher  *BusinessMatcher
	scorer   *LeadScorer
	dedup    *Deduplicator
}

func NewAnalyst(store *memory.Store) *Analyst {
	return &Analyst{
		Store:    store,
		verifier: NewEmailVerifier(),
		scanner:  NewWebsiteScanner(),
		matcher:  NewBusinessMatcher(),
		scorer:   NewLeadScorer(),
		dedup:    NewDeduplicator(store),
	}
}

func (a *Analyst) Analyze(lead *memory.Lead) *memory.Lead {
	if lead.Email != "" {
		if a.dedup.IsDuplicate(lead.Email) {
			lead.Status = "duplicate"
			return lead
		}
	}

	if lead.Email != "" {
		valid, err := a.verifier.Verify(lead.Email)
		if err != nil {
			log.Printf("email verify error: %v", err)
		}
		if !valid {
			lead.Status = "invalid_email"
			return lead
		}
	}

	if lead.Website != "" {
		scanResult, err := a.scanner.Scan(lead.Website)
		if err != nil {
			log.Printf("website scan error: %v", err)
		} else {
			if lead.Company == "" {
				lead.Company = scanResult.Title
			}
			if lead.Notes == "" {
				lead.Notes = scanResult.Description
			}
		}
	}

	matchScore := a.matcher.Match(lead, "FMCG Foods, Superfoods, Nutraceuticals, Horticulture, Eco-Packaging")
	lead.Score = a.scorer.Score(lead, matchScore)

	if lead.Email != "" {
		lead.Status = "verified"
	}

	return lead
}
