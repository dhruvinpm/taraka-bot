package telegram

import (
	telebot "gopkg.in/telebot.v4"
)

func ReplyApprovalKeyboard() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}
	btnSend := menu.Data("✅ Send", "reply_send")
	btnEdit := menu.Data("✏️ Edit", "reply_edit")
	btnSkip := menu.Data("❌ Skip", "reply_skip")
	menu.Inline(menu.Row(btnSend, btnEdit, btnSkip))
	return menu
}
