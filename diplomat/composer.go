package diplomat

import (
	"context"

	"github.com/dhruvinpm/taraka-bot/llm"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type Composer struct {
	llm       llm.Provider
	humanizer *Humanizer
}

func NewComposer(llmProvider llm.Provider) *Composer {
	return &Composer{
		llm:       llmProvider,
		humanizer: NewHumanizer(),
	}
}

func (c *Composer) ComposeInitial(ctx context.Context, lead *memory.Lead) (string, string, error) {
	prompt := llm.InitialOutreachPrompt(lead.Name, lead.Company, lead.BusinessType, lead.Country, lead.Website, lead.Notes)
	body, err := c.llm.Generate(ctx, prompt)
	if err != nil {
		return "", "", err
	}
	body = c.humanizer.Humanize(body)
	subject := "Business Partnership Opportunity - Taraka International"
	return subject, body, nil
}

func (c *Composer) ComposeFollowUp(ctx context.Context, lead *memory.Lead, followUpNum int) (string, string, error) {
	var prompt string
	switch followUpNum {
	case 1:
		prompt = llm.FollowUp1Prompt(lead.Name, lead.Company)
	case 2:
		prompt = llm.FollowUp2Prompt(lead.Name, lead.Company)
	default:
		prompt = llm.FollowUp3Prompt(lead.Name, lead.Company)
	}
	body, err := c.llm.Generate(ctx, prompt)
	if err != nil {
		return "", "", err
	}
	body = c.humanizer.Humanize(body)
	subject := "Following up - Taraka International"
	return subject, body, nil
}

func (c *Composer) ComposeReply(ctx context.Context, lead *memory.Lead, replyText string) (string, string, error) {
	prompt := llm.ReplyResponsePrompt(replyText)
	body, err := c.llm.Generate(ctx, prompt)
	if err != nil {
		return "", "", err
	}
	body = c.humanizer.Humanize(body)
	subject := "Re: Business Partnership - Taraka International"
	return subject, body, nil
}
