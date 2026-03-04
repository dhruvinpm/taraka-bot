package diplomat

import (
	"context"
	"log"
	"sync"
	"time"

	localmail "github.com/dhruvinpm/taraka-bot/gmail"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type EmailSender struct {
	sender    *localmail.Sender
	warmup    *WarmupManager
	mu        sync.Mutex
	sentToday int
	lastReset time.Time
}

func NewEmailSender(sender *localmail.Sender, warmup *WarmupManager) *EmailSender {
	return &EmailSender{
		sender:    sender,
		warmup:    warmup,
		lastReset: time.Now(),
	}
}

func (e *EmailSender) Send(ctx context.Context, lead *memory.Lead, subject, body string) error {
	e.mu.Lock()
	if time.Since(e.lastReset) > 24*time.Hour {
		e.sentToday = 0
		e.lastReset = time.Now()
	}
	limit := e.warmup.GetDailyLimit()
	if e.sentToday >= limit {
		e.mu.Unlock()
		return nil
	}
	e.sentToday++
	e.mu.Unlock()

	if err := e.sender.Send(lead.Email, subject, body); err != nil {
		return err
	}
	log.Printf("sent email to %s (%s)", lead.Email, lead.Company)
	return nil
}
