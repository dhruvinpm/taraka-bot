package hunter

import (
	"regexp"
	"strings"

	"github.com/dhruvinpm/taraka-bot/memory"
	"github.com/gocolly/colly/v2"
)

type WebsiteHunter struct{}

func NewWebsiteHunter() *WebsiteHunter {
	return &WebsiteHunter{}
}

var emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
var phoneRegex = regexp.MustCompile(`[\+]?[(]?[0-9]{1,4}[)]?[-\s\./0-9]{8,15}`)

func (w *WebsiteHunter) ScrapeWebsite(siteURL string) (*memory.Lead, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.MaxDepth(2),
	)

	lead := &memory.Lead{
		Website: siteURL,
		Source:  "website_scrape",
		Status:  "raw",
	}

	c.OnHTML("title", func(e *colly.HTMLElement) {
		if lead.Company == "" {
			lead.Company = strings.TrimSpace(e.Text)
		}
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		text := e.Text
		emails := emailRegex.FindAllString(text, -1)
		if len(emails) > 0 && lead.Email == "" {
			lead.Email = emails[0]
		}
		phones := phoneRegex.FindAllString(text, -1)
		if len(phones) > 0 && lead.Phone == "" {
			lead.Phone = phones[0]
		}
	})

	if err := c.Visit(siteURL); err != nil {
		return nil, err
	}

	return lead, nil
}
