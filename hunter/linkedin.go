package hunter

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dhruvinpm/taraka-bot/browser"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type LinkedInHunter struct {
	pool *browser.BrowserPool
}

func NewLinkedInHunter(pool *browser.BrowserPool) *LinkedInHunter {
	return &LinkedInHunter{pool: pool}
}

func (l *LinkedInHunter) SearchLinkedIn(query string) ([]*memory.Lead, error) {
	if l.pool == nil {
		return nil, fmt.Errorf("browser pool not available")
	}

	b, err := l.pool.Get()
	if err != nil {
		return nil, fmt.Errorf("get browser: %w", err)
	}
	defer l.pool.Release(b)

	page, err := browser.StealthPage(b)
	if err != nil {
		return nil, fmt.Errorf("stealth page: %w", err)
	}
	defer page.MustClose()

	// Google dork for LinkedIn company pages — no login required.
	dorkQuery := fmt.Sprintf("site:linkedin.com/company %s", query)
	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s&num=20&hl=en", url.QueryEscape(dorkQuery))

	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("navigate: %w", err)
	}
	if err := page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("wait load: %w", err)
	}

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
		if err != nil || href == nil || !strings.Contains(*href, "linkedin.com") {
			continue
		}

		snippet := ""
		if snippetEl, err := el.Element("div.VwiC3b"); err == nil && snippetEl != nil {
			if s, err := snippetEl.Text(); err == nil {
				snippet = s
			}
		}

		leads = append(leads, &memory.Lead{
			Name:   title,
			Website: *href,
			Source:  "linkedin",
			Notes:   snippet,
			Status:  "raw",
		})
	}

	return leads, nil
}
