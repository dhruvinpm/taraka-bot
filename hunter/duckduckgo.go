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
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		colly.MaxDepth(1),
	)

	var leads []*memory.Lead

	// DDG Lite uses a table-based layout. Each result title and URL is in
	// <a class="result-link"> and the snippet is in <td class="result-snippet">.
	c.OnHTML("a.result-link", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.Text)
		link := e.Attr("href")
		if title == "" {
			return
		}
		if strings.HasPrefix(link, "//") {
			link = "https:" + link
		} else if !strings.HasPrefix(link, "http") {
			return
		}

		lead := &memory.Lead{
			Company: title,
			Website: link,
			Source:  "duckduckgo",
			Status:  "raw",
		}
		leads = append(leads, lead)
	})

	// Attach snippets to the most recently parsed lead.
	c.OnHTML("td.result-snippet", func(e *colly.HTMLElement) {
		snippet := strings.TrimSpace(e.Text)
		if len(leads) > 0 && leads[len(leads)-1].Notes == "" {
			leads[len(leads)-1].Notes = snippet
		}
	})

	// DDG Lite returns simple static HTML that Colly can parse without JS.
	searchURL := fmt.Sprintf("https://lite.duckduckgo.com/lite/?q=%s", url.QueryEscape(query))
	if err := c.Visit(searchURL); err != nil {
		return nil, err
	}

	return leads, nil
}
