package hunter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dhruvinpm/taraka-bot/memory"
)

type RedditHunter struct{}

func NewRedditHunter() *RedditHunter {
	return &RedditHunter{}
}

type redditSearchResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				Title     string `json:"title"`
				URL       string `json:"url"`
				Selftext  string `json:"selftext"`
				Subreddit string `json:"subreddit"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func (r *RedditHunter) SearchReddit(niche string) ([]*memory.Lead, error) {
	searchQuery := niche + " importer buyer"
	apiURL := fmt.Sprintf("https://www.reddit.com/search.json?q=%s&sort=new&limit=25", url.QueryEscape(searchQuery))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	// Reddit requires a descriptive User-Agent for API access.
	req.Header.Set("User-Agent", "TarakaBot/1.0 lead-research-tool")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reddit API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result redditSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var leads []*memory.Lead
	for _, child := range result.Data.Children {
		post := child.Data
		if post.Title == "" {
			continue
		}

		link := post.URL
		if !strings.HasPrefix(link, "http") {
			link = "https://www.reddit.com" + link
		}

		notes := post.Selftext
		if len(notes) > 500 {
			notes = notes[:500]
		}

		leads = append(leads, &memory.Lead{
			Company: post.Title,
			Website: link,
			Source:  "reddit",
			Notes:   notes,
			Status:  "raw",
		})
	}

	return leads, nil
}
