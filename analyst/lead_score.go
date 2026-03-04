package analyst

import "github.com/dhruvinpm/taraka-bot/memory"

type LeadScorer struct{}

func NewLeadScorer() *LeadScorer {
	return &LeadScorer{}
}

func (ls *LeadScorer) Score(lead *memory.Lead, matchScore int) int {
	score := matchScore

	if lead.Email != "" {
		score += 20
	}
	if lead.Website != "" {
		score += 10
	}
	if lead.Phone != "" {
		score += 10
	}
	if lead.Company != "" {
		score += 5
	}
	if lead.Country != "" {
		score += 5
	}

	if score > 100 {
		score = 100
	}
	return score
}
