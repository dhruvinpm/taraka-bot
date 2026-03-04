package telegram

import (
	"fmt"
	"log"

	"github.com/dhruvinpm/taraka-bot/memory"
	telebot "gopkg.in/telebot.v4"
)

type Notifier struct {
	bot         *telebot.Bot
	adminChatID int64
}

func NewNotifier(bot *telebot.Bot, adminChatID int64) *Notifier {
	return &Notifier{bot: bot, adminChatID: adminChatID}
}

func (n *Notifier) SendNotification(text string) {
	if n.bot == nil {
		return
	}
	if _, err := n.bot.Send(telebot.ChatID(n.adminChatID), text); err != nil {
		log.Printf("notification send error: %v", err)
	}
}

func (n *Notifier) SendDailySummary(stats map[string]int) {
	msg := "📊 *Daily Summary*\n"
	for status, count := range stats {
		msg += fmt.Sprintf("• %s: %d\n", status, count)
	}
	n.SendNotification(msg)
}

func (n *Notifier) SendReplyAlert(lead *memory.Lead, replyText string) {
	msg := fmt.Sprintf("📬 *Reply Received*\nFrom: %s (%s)\nEmail: %s\n\n%s",
		lead.Name, lead.Company, lead.Email, replyText)
	n.SendNotification(msg)
}
