package diplomat

import (
	"regexp"
	"strings"
)

var aiPhrases = []string{
	"I hope this email finds you well",
	"I hope this finds you well",
	"Please don't hesitate to",
	"Please feel free to",
	"I'd be delighted to",
	"I'd be happy to",
	"Looking forward to hearing from you",
}

var markdownBold = regexp.MustCompile(`\*\*(.+?)\*\*`)
var markdownItalic = regexp.MustCompile(`\*(.+?)\*`)
var markdownHeader = regexp.MustCompile(`(?m)^#+\s+`)
var bulletPoint = regexp.MustCompile(`(?m)^[\-\*]\s+`)
var numberedList = regexp.MustCompile(`(?m)^\d+\.\s+`)
var emDashPattern = regexp.MustCompile(`—\s*`)

type Humanizer struct{}

func NewHumanizer() *Humanizer {
	return &Humanizer{}
}

func (h *Humanizer) Humanize(text string) string {
	for _, phrase := range aiPhrases {
		text = strings.ReplaceAll(text, phrase, "")
	}

	text = markdownBold.ReplaceAllString(text, "$1")
	text = markdownItalic.ReplaceAllString(text, "$1")
	text = markdownHeader.ReplaceAllString(text, "")
	text = bulletPoint.ReplaceAllString(text, "")
	text = numberedList.ReplaceAllString(text, "")
	text = emDashPattern.ReplaceAllString(text, "- ")

	lines := strings.Split(text, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, "\n")
}
