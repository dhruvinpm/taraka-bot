package diplomat

import (
	"context"
	"log"

	localmail "github.com/dhruvinpm/taraka-bot/gmail"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type ReplyHandler struct {
	client   *localmail.Client
	store    *memory.Store
	composer *Composer
	sender   *EmailSender
}

func NewReplyHandler(client *localmail.Client, store *memory.Store, composer *Composer, sender *EmailSender) *ReplyHandler {
	return &ReplyHandler{
		client:   client,
		store:    store,
		composer: composer,
		sender:   sender,
	}
}

func (r *ReplyHandler) CheckReplies(ctx context.Context) error {
	messages, err := r.client.ListMessages("is:unread in:inbox", 50)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		fullMsg, err := r.client.GetMessage(msg.Id)
		if err != nil {
			log.Printf("get message error: %v", err)
			continue
		}

		threadID := fullMsg.ThreadId
		lead, err := r.store.GetLeadByThreadID(threadID)
		if err != nil {
			continue
		}

		body := localmail.GetMessageBody(fullMsg)
		if body != "" {
			r.ProcessReply(ctx, lead, body)
		}

		r.client.MarkAsRead(msg.Id)
	}

	return nil
}

func (r *ReplyHandler) ProcessReply(ctx context.Context, lead *memory.Lead, replyText string) {
	if err := r.store.MarkReplyReceived(lead.ID); err != nil {
		log.Printf("mark reply error: %v", err)
		return
	}

	subject, body, err := r.composer.ComposeReply(ctx, lead, replyText)
	if err != nil {
		log.Printf("compose reply error: %v", err)
		return
	}

	if err := r.sender.Send(ctx, lead, subject, body); err != nil {
		log.Printf("send reply error: %v", err)
	}
}
