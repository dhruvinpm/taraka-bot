package hunter

import (
	"fmt"
	"net/url"

	"github.com/dhruvinpm/taraka-bot/memory"
	"github.com/gocolly/colly/v2"
)

type DirectoryHunter struct{}

func NewDirectoryHunter() *DirectoryHunter {
	return &DirectoryHunter{}
}

func (d *DirectoryHunter) SearchDirectories(query, country string) ([]*memory.Lead, error) {
	var all []*memory.Lead

	kompassLeads, _ := d.searchKompass(query, country)
	all = append(all, kompassLeads...)

	euroLeads, _ := d.searchEuropages(query, country)
	all = append(all, euroLeads...)

	return all, nil
}

func (d *DirectoryHunter) searchKompass(query, country string) ([]*memory.Lead, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.MaxDepth(1),
	)

	var leads []*memory.Lead

	c.OnHTML(".companyCard", func(e *colly.HTMLElement) {
		name := e.ChildText(".companyCard-title")
		website := e.ChildAttr("a", "href")
		if name == "" {
			return
		}
		leads = append(leads, &memory.Lead{
			Company: name,
			Website: website,
			Country: country,
			Source:  "kompass",
			Status:  "raw",
		})
	})

	searchURL := fmt.Sprintf("https://www.kompass.com/search/?text=%s", url.QueryEscape(query))
	c.Visit(searchURL)

	return leads, nil
}

func (d *DirectoryHunter) searchEuropages(query, country string) ([]*memory.Lead, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.MaxDepth(1),
	)

	var leads []*memory.Lead

	c.OnHTML(".company-card", func(e *colly.HTMLElement) {
		name := e.ChildText(".company-name")
		website := e.ChildAttr("a", "href")
		if name == "" {
			return
		}
		leads = append(leads, &memory.Lead{
			Company: name,
			Website: website,
			Country: country,
			Source:  "europages",
			Status:  "raw",
		})
	})

	searchURL := fmt.Sprintf("https://www.europages.co.uk/en/search?q=%s", url.QueryEscape(query))
	c.Visit(searchURL)

	return leads, nil
}
