package analyst

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

type ScanResult struct {
	Title       string
	Description string
	Keywords    string
	ContactInfo string
}

type WebsiteScanner struct{}

func NewWebsiteScanner() *WebsiteScanner {
	return &WebsiteScanner{}
}

func (w *WebsiteScanner) Scan(siteURL string) (*ScanResult, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.MaxDepth(1),
	)

	result := &ScanResult{}

	c.OnHTML("title", func(e *colly.HTMLElement) {
		result.Title = strings.TrimSpace(e.Text)
	})

	c.OnHTML("meta[name='description']", func(e *colly.HTMLElement) {
		result.Description = e.Attr("content")
	})

	c.OnHTML("meta[name='keywords']", func(e *colly.HTMLElement) {
		result.Keywords = e.Attr("content")
	})

	if err := c.Visit(siteURL); err != nil {
		return result, err
	}

	return result, nil
}
