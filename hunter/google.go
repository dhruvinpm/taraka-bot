package hunter

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dhruvinpm/taraka-bot/browser"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type GoogleHunter struct {
	pool *browser.BrowserPool
}

func NewGoogleHunter(pool *browser.BrowserPool) *GoogleHunter {
	return &GoogleHunter{pool: pool}
}

func (g *GoogleHunter) SearchGoogle(query, country string) ([]*memory.Lead, error) {
	if g.pool == nil {
		return nil, fmt.Errorf("browser pool not available")
	}

	b, err := g.pool.Get()
	if err != nil {
		return nil, fmt.Errorf("get browser: %w", err)
	}
	defer g.pool.Release(b)

	page, err := browser.StealthPage(b)
	if err != nil {
		return nil, fmt.Errorf("stealth page: %w", err)
	}
	defer page.MustClose()

	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s&num=30&hl=en", url.QueryEscape(query))
	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("navigate: %w", err)
	}
	if err := page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("wait load: %w", err)
	}

	// Brief human-like pause to let results render.
	time.Sleep(2 * time.Second)

	elements, err := page.Elements("div.g")
	if err != nil {
		return nil, fmt.Errorf("get elements: %w", err)
	}

	var leads []*memory.Lead
	for _, el := range elements {
		titleEl, err := el.Element("h3")
		if err != nil || titleEl == nil {
			continue
		}
		title, err := titleEl.Text()
		if err != nil || title == "" {
			continue
		}

		linkEl, err := el.Element("a")
		if err != nil || linkEl == nil {
			continue
		}
		href, err := linkEl.Attribute("href")
		if err != nil || href == nil || !strings.HasPrefix(*href, "http") {
			continue
		}

		snippet := ""
		if snippetEl, err := el.Element("div.VwiC3b"); err == nil && snippetEl != nil {
			if s, err := snippetEl.Text(); err == nil {
				snippet = s
			}
		}

		leads = append(leads, &memory.Lead{
			Company: title,
			Website: *href,
			Country: country,
			Source:  "google",
			Notes:   snippet,
			Status:  "raw",
		})
	}

	return leads, nil
}
