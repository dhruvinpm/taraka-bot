package hunter

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dhruvinpm/taraka-bot/memory"
	"github.com/gocolly/colly/v2"
)

type RedditHunter struct{}

func NewRedditHunter() *RedditHunter {
	return &RedditHunter{}
}

func (r *RedditHunter) SearchReddit(niche string) ([]*memory.Lead, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.MaxDepth(1),
	)

	var leads []*memory.Lead

	c.OnHTML(".thing", func(e *colly.HTMLElement) {
		title := e.ChildText("a.title")
		link := e.ChildAttr("a.title", "href")

		if title == "" {
			return
		}
		if !strings.HasPrefix(link, "http") {
			link = "https://www.reddit.com" + link
		}

		lead := &memory.Lead{
			Company: title,
			Website: link,
			Source:  "reddit",
			Status:  "raw",
		}
		leads = append(leads, lead)
	})

	searchURL := fmt.Sprintf("https://old.reddit.com/search?q=%s&sort=new", url.QueryEscape(niche+" importer buyer"))
	if err := c.Visit(searchURL); err != nil {
		return nil, err
	}

	return leads, nil
}
