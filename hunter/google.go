package hunter

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dhruvinpm/taraka-bot/memory"
	"github.com/gocolly/colly/v2"
)

type GoogleHunter struct{}

func NewGoogleHunter() *GoogleHunter {
	return &GoogleHunter{}
}

func (g *GoogleHunter) SearchGoogle(query, country string) ([]*memory.Lead, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.MaxDepth(1),
	)

	var leads []*memory.Lead

	c.OnHTML("div.g", func(e *colly.HTMLElement) {
		title := e.ChildText("h3")
		link := e.ChildAttr("a", "href")
		snippet := e.ChildText("div.VwiC3b")

		if title == "" || link == "" {
			return
		}
		if !strings.HasPrefix(link, "http") {
			return
		}

		lead := &memory.Lead{
			Company: title,
			Website: link,
			Country: country,
			Source:  "google",
			Notes:   snippet,
			Status:  "raw",
		}
		leads = append(leads, lead)
	})

	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s&num=30", url.QueryEscape(query))
	if err := c.Visit(searchURL); err != nil {
		return nil, err
	}

	return leads, nil
}
