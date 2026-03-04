package hunter

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dhruvinpm/taraka-bot/memory"
	"github.com/gocolly/colly/v2"
)

type DuckDuckGoHunter struct{}

func NewDuckDuckGoHunter() *DuckDuckGoHunter {
	return &DuckDuckGoHunter{}
}

func (d *DuckDuckGoHunter) SearchDDG(query string) ([]*memory.Lead, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.MaxDepth(1),
	)

	var leads []*memory.Lead

	c.OnHTML(".result", func(e *colly.HTMLElement) {
		title := e.ChildText(".result__title")
		link := e.ChildAttr(".result__url", "href")
		snippet := e.ChildText(".result__snippet")

		if title == "" {
			return
		}
		if !strings.HasPrefix(link, "http") {
			link = "https://" + link
		}

		lead := &memory.Lead{
			Company: title,
			Website: link,
			Source:  "duckduckgo",
			Notes:   snippet,
			Status:  "raw",
		}
		leads = append(leads, lead)
	})

	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))
	if err := c.Visit(searchURL); err != nil {
		return nil, err
	}

	return leads, nil
}
