package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

func StealthPage(browser *rod.Browser) (*rod.Page, error) {
	page, err := stealth.Page(browser)
	if err != nil {
		return nil, err
	}
	return page, nil
}
