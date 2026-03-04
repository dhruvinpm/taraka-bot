package diplomat

import (
	"context"
	"log"
	"time"

	"github.com/dhruvinpm/taraka-bot/memory"
)

type FollowUpEngine struct {
	store    *memory.Store
	composer *Composer
	sender   *EmailSender
}

func NewFollowUpEngine(store *memory.Store, composer *Composer, sender *EmailSender) *FollowUpEngine {
	return &FollowUpEngine{store: store, composer: composer, sender: sender}
}

var followUpDays = []int{5, 10, 15}
var followUpStatuses = []string{"follow_up_1", "follow_up_2", "follow_up_3"}

func (f *FollowUpEngine) ProcessFollowUps(ctx context.Context) error {
	leads, err := f.store.GetDueFollowUps()
	if err != nil {
		return err
	}

	for _, lead := range leads {
		followUpNum := lead.FollowUpsSent + 1
		if followUpNum > 3 {
			f.store.UpdateLeadStatus(lead.ID, "exhausted")
			continue
		}

		subject, body, err := f.composer.ComposeFollowUp(ctx, lead, followUpNum)
		if err != nil {
			log.Printf("compose follow-up error: %v", err)
			continue
		}

		if err := f.sender.Send(ctx, lead, subject, body); err != nil {
			log.Printf("send follow-up error: %v", err)
			continue
		}

		nextFollowUp := time.Now().AddDate(0, 0, followUpDays[followUpNum-1])
		status := followUpStatuses[followUpNum-1]
		if err := f.store.UpdateFollowUp(lead.ID, nextFollowUp, status); err != nil {
			log.Printf("update follow-up error: %v", err)
		}
	}

	return nil
}

func (f *FollowUpEngine) ScheduleFollowUp(lead *memory.Lead, followUpNum int) error {
	if followUpNum < 1 || followUpNum > len(followUpDays) {
		return nil
	}
	nextFollowUp := time.Now().AddDate(0, 0, followUpDays[followUpNum-1])
	return f.store.UpdateFollowUp(lead.ID, nextFollowUp, "emailed")
}
