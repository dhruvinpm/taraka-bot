package hunter

import (
	"fmt"
	"net/url"

	"github.com/dhruvinpm/taraka-bot/memory"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type LinkedInHunter struct {
	browser *rod.Browser
}

func NewLinkedInHunter(browser *rod.Browser) *LinkedInHunter {
	return &LinkedInHunter{browser: browser}
}

func (l *LinkedInHunter) SearchLinkedIn(query string) ([]*memory.Lead, error) {
	dorkQuery := fmt.Sprintf("site:linkedin.com/in/ %s", query)
	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s", url.QueryEscape(dorkQuery))

	page, err := l.browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return nil, fmt.Errorf("create page: %w", err)
	}
	defer page.MustClose()

	if err := page.Navigate(searchURL); err != nil {
		return nil, err
	}
	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	elements, err := page.Elements("div.g h3")
	if err != nil {
		return nil, err
	}

	var leads []*memory.Lead
	for _, el := range elements {
		text, err := el.Text()
		if err != nil {
			continue
		}
		lead := &memory.Lead{
			Name:   text,
			Source: "linkedin",
			Status: "raw",
		}
		leads = append(leads, lead)
	}

	return leads, nil
}
