package diplomat

import (
	"context"
	"log"

	localmail "github.com/dhruvinpm/taraka-bot/gmail"
	"github.com/dhruvinpm/taraka-bot/llm"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type Diplomat struct {
	store        *memory.Store
	composer     *Composer
	sender       *EmailSender
	followup     *FollowUpEngine
	replyHandler *ReplyHandler
	tracker      *GmailTracker
}

type Config struct {
	SMTPConfig  localmail.SMTPConfig
	WarmupStart int
	WarmupInc   int
	WarmupMax   int
}

func NewDiplomat(store *memory.Store, gmailClient *localmail.Client, cfg Config, llmProvider llm.Provider) *Diplomat {
	warmup := NewWarmupManager(cfg.WarmupStart, cfg.WarmupInc, cfg.WarmupMax)
	gmailSender := localmail.NewSender(cfg.SMTPConfig)
	sender := NewEmailSender(gmailSender, warmup)
	composer := NewComposer(llmProvider)
	followup := NewFollowUpEngine(store, composer, sender)
	replyHandler := NewReplyHandler(gmailClient, store, composer, sender)
	tracker := NewGmailTracker(gmailClient, store)

	return &Diplomat{
		store:        store,
		composer:     composer,
		sender:       sender,
		followup:     followup,
		replyHandler: replyHandler,
		tracker:      tracker,
	}
}

func (d *Diplomat) SendInitialOutreach(ctx context.Context, lead *memory.Lead) error {
	subject, body, err := d.composer.ComposeInitial(ctx, lead)
	if err != nil {
		return err
	}

	if err := d.sender.Send(ctx, lead, subject, body); err != nil {
		return err
	}

	if err := d.store.UpdateLeadStatus(lead.ID, "emailed"); err != nil {
		log.Printf("update status error: %v", err)
	}

	return nil
}

func (d *Diplomat) ProcessFollowUps(ctx context.Context) error {
	return d.followup.ProcessFollowUps(ctx)
}

func (d *Diplomat) CheckReplies(ctx context.Context) error {
	return d.replyHandler.CheckReplies(ctx)
}
