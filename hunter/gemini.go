package hunter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dhruvinpm/taraka-bot/llm"
	"github.com/dhruvinpm/taraka-bot/memory"
)

// GeminiHunter uses the Gemini LLM to find business leads.
// This is the primary hunter source because direct web scraping
// (Google, DuckDuckGo) is blocked by anti-bot protections.
type GeminiHunter struct {
	llm llm.Provider
}

func NewGeminiHunter(llmProvider llm.Provider) *GeminiHunter {
	return &GeminiHunter{llm: llmProvider}
}

type geminiLead struct {
	CompanyName         string `json:"company_name"`
	ContactName         string `json:"contact_name"`
	Email               string `json:"email"`
	Phone               string `json:"phone"`
	Website             string `json:"website"`
	BusinessDescription string `json:"business_description"`
	Country             string `json:"country"`
}

// SearchGemini asks Gemini to find business leads matching the query and country,
// returning them as structured Lead objects. Gemini can access the web and return
// real company information, making this far more reliable than scraping.
// maxLeads controls how many companies are requested (0 means use the default of 20).
func (g *GeminiHunter) SearchGemini(ctx context.Context, query, country string, maxLeads int) ([]*memory.Lead, error) {
	if g.llm == nil {
		return nil, fmt.Errorf("no LLM provider configured")
	}
	if maxLeads <= 0 {
		maxLeads = 20
	}

	prompt := fmt.Sprintf(`You are a business lead finder for Taraka International, an export/import company specialising in FMCG Foods, Superfoods, Nutraceuticals, Horticulture, and Eco-Packaging.

Find %d real companies matching this search: "%s" in %s.
Focus on importers, distributors, buyers, wholesalers, supermarkets, health food stores, or retailers.

Return ONLY a valid JSON array (no markdown fences, no explanation) with this exact structure:
[
  {
    "company_name": "Company Name",
    "contact_name": "Contact Person (if known, else empty string)",
    "email": "contact@company.com (valid email only, else empty string)",
    "phone": "+44... (if known, else empty string)",
    "website": "https://www.company.com (if known, else empty string)",
    "business_description": "Brief description of what they do",
    "country": "%s"
  }
]

Rules:
- Only include companies that could realistically be leads for an exporter.
- Prefer entries where you know the email address.
- Do NOT invent email addresses — use empty string "" if unknown.
- Return valid JSON only.`, maxLeads, query, country, country)

	resp, err := g.llm.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("gemini search: %w", err)
	}

	resp = extractJSONArray(resp)

	var geminiLeads []geminiLead
	if err := json.Unmarshal([]byte(resp), &geminiLeads); err != nil {
		return nil, fmt.Errorf("parse gemini response: %w (raw: %.200s)", err, resp)
	}

	var leads []*memory.Lead
	for _, gl := range geminiLeads {
		if gl.CompanyName == "" {
			continue
		}
		leads = append(leads, &memory.Lead{
			Company:      gl.CompanyName,
			Name:         gl.ContactName,
			Email:        gl.Email,
			Phone:        gl.Phone,
			Website:      gl.Website,
			Country:      country,
			Source:       "gemini",
			Notes:        gl.BusinessDescription,
			Status:       "raw",
		})
	}
	return leads, nil
}

// extractJSONArray strips markdown code fences and extracts the JSON array portion
// from a Gemini response that may contain extra prose.
func extractJSONArray(s string) string {
	// Strip markdown code fences (```json ... ```)
	s = strings.ReplaceAll(s, "```json", "")
	s = strings.ReplaceAll(s, "```", "")
	s = strings.TrimSpace(s)

	start := strings.Index(s, "[")
	end := strings.LastIndex(s, "]")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}
