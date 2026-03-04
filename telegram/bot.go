package telegram

import (
	"log"
	"time"

	"github.com/dhruvinpm/taraka-bot/core"
	"github.com/dhruvinpm/taraka-bot/memory"
	telebot "gopkg.in/telebot.v4"
)

type TelegramBot struct {
	bot      *telebot.Bot
	engine   *core.Engine
	store    *memory.Store
	handlers *Handlers
	notifier *Notifier
}

func New(token string, adminChatID int64, engine *core.Engine, store *memory.Store) (*TelegramBot, error) {
	settings := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	h := NewHandlers(engine, store)
	n := NewNotifier(bot, adminChatID)

	tb := &TelegramBot{
		bot:      bot,
		engine:   engine,
		store:    store,
		handlers: h,
		notifier: n,
	}

	tb.registerHandlers()
	return tb, nil
}

func (tb *TelegramBot) registerHandlers() {
	tb.bot.Handle("/start", tb.handlers.HandleStart)
	tb.bot.Handle("/hunt", tb.handlers.HandleHunt)
	tb.bot.Handle("/stats", tb.handlers.HandleStats)
	tb.bot.Handle("/pause", tb.handlers.HandlePause)
	tb.bot.Handle("/resume", tb.handlers.HandleResume)
	tb.bot.Handle("/leads", tb.handlers.HandleLeads)
	tb.bot.Handle("/query", tb.handlers.HandleQuery)
	tb.bot.Handle("/mode", tb.handlers.HandleMode)
	tb.bot.Handle("/export", tb.handlers.HandleExport)
	tb.bot.Handle("/status", tb.handlers.HandleStatus)
	tb.bot.Handle("/follow", tb.handlers.HandleFollow)
	tb.bot.Handle("/stop", tb.handlers.HandleStop)
	tb.bot.Handle("/warmup", tb.handlers.HandleWarmup)
}

func (tb *TelegramBot) Start() {
	log.Println("telegram bot starting...")
	tb.bot.Start()
}

func (tb *TelegramBot) Stop() {
	tb.bot.Stop()
}

func (tb *TelegramBot) GetNotifier() *Notifier {
	return tb.notifier
}
