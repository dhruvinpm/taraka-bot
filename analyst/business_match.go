package analyst

import (
	"strings"

	"github.com/dhruvinpm/taraka-bot/memory"
)

var targetKeywords = []string{
	"import", "importer", "buyer", "distributor", "wholesale", "retailer",
	"food", "grocery", "supermarket", "organic", "health", "nutraceutical",
	"packaging", "eco", "sustainable", "fmcg", "consumer goods",
	"horticulture", "agriculture", "cocopeat",
}

var targetMarkets = map[string]bool{
	"UAE": true, "UK": true, "USA": true, "Australia": true,
	"Singapore": true, "Canada": true, "New Zealand": true,
	"Japan": true, "South Korea": true,
}

type BusinessMatcher struct{}

func NewBusinessMatcher() *BusinessMatcher {
	return &BusinessMatcher{}
}

func (b *BusinessMatcher) Match(lead *memory.Lead, companyProfile string) int {
	score := 0
	text := strings.ToLower(lead.Company + " " + lead.Notes + " " + lead.BusinessType)

	for _, kw := range targetKeywords {
		if strings.Contains(text, kw) {
			score += 5
		}
	}

	if targetMarkets[lead.Country] {
		score += 20
	}

	if score > 100 {
		score = 100
	}
	return score
}
