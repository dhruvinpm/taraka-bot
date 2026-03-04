package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/dhruvinpm/taraka-bot/core"
	"github.com/dhruvinpm/taraka-bot/memory"
	telebot "gopkg.in/telebot.v4"
)

type Handlers struct {
	engine *core.Engine
	store  *memory.Store
}

func NewHandlers(engine *core.Engine, store *memory.Store) *Handlers {
	return &Handlers{engine: engine, store: store}
}

func (h *Handlers) HandleStart(c telebot.Context) error {
	return c.Send("👋 TarakaBot is online!\n\nCommands:\n/hunt - Start lead hunting\n/stats - Show statistics\n/pause - Pause operations\n/resume - Resume operations\n/leads - List recent leads\n/query - Query database\n/status - Show status")
}

func (h *Handlers) HandleHunt(c telebot.Context) error {
	args := c.Args()
	country := ""
	niche := ""
	if len(args) >= 1 {
		country = args[0]
	}
	if len(args) >= 2 {
		niche = strings.Join(args[1:], " ")
	}

	c.Send("🔍 Starting hunt...")
	go func() {
		ctx := context.Background()
		notify := func(msg string) {
			c.Send(msg)
		}
		if err := h.engine.RunHunt(ctx, country, niche, notify); err != nil {
			c.Send(fmt.Sprintf("❌ Hunt error: %v", err))
		}
	}()
	return nil
}

func (h *Handlers) HandleStats(c telebot.Context) error {
	stats, err := h.store.GetStats()
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Error: %v", err))
	}
	msg := "📊 Lead Statistics:\n"
	total := 0
	for status, count := range stats {
		msg += fmt.Sprintf("• %s: %d\n", status, count)
		total += count
	}
	msg += fmt.Sprintf("\nTotal: %d", total)
	return c.Send(msg)
}

func (h *Handlers) HandlePause(c telebot.Context) error {
	h.engine.Pause()
	return c.Send("⏸ Operations paused")
}

func (h *Handlers) HandleResume(c telebot.Context) error {
	h.engine.Resume()
	return c.Send("▶️ Operations resumed")
}

func (h *Handlers) HandleLeads(c telebot.Context) error {
	leads, err := h.store.GetLeadsByStatus("verified", 10)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Error: %v", err))
	}
	if len(leads) == 0 {
		return c.Send("No verified leads found")
	}
	msg := "📋 Recent Verified Leads:\n"
	for i, lead := range leads {
		msg += fmt.Sprintf("%d. %s (%s) - Score: %d\n", i+1, lead.Company, lead.Country, lead.Score)
	}
	return c.Send(msg)
}

func (h *Handlers) HandleQuery(c telebot.Context) error {
	query := strings.Join(c.Args(), " ")
	if query == "" {
		return c.Send("Usage: /query <SQL or natural language>")
	}
	return c.Send(fmt.Sprintf("Query received: %s\n(Use /query with SQL for direct queries)", query))
}

func (h *Handlers) HandleMode(c telebot.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Usage: /mode <gentle|moderate|aggressive>")
	}
	return c.Send(fmt.Sprintf("Mode set to: %s", args[0]))
}

func (h *Handlers) HandleExport(c telebot.Context) error {
	return c.Send("Export feature coming soon")
}

func (h *Handlers) HandleStatus(c telebot.Context) error {
	status := "▶️ Running"
	if h.engine.IsPaused() {
		status = "⏸ Paused"
	}
	count, _ := h.store.CountLeads()
	return c.Send(fmt.Sprintf("Status: %s\nTotal leads: %d", status, count))
}

func (h *Handlers) HandleFollow(c telebot.Context) error {
	go func() {
		ctx := context.Background()
		if err := h.engine.RunOutreach(ctx); err != nil {
			c.Send(fmt.Sprintf("❌ Follow-up error: %v", err))
		} else {
			c.Send("✅ Follow-ups processed")
		}
	}()
	return c.Send("📨 Processing follow-ups...")
}

func (h *Handlers) HandleStop(c telebot.Context) error {
	h.engine.Pause()
	return c.Send("🛑 Bot operations stopped")
}

func (h *Handlers) HandleWarmup(c telebot.Context) error {
	return c.Send("🌡 Warmup status: Active\nCheck logs for current daily limit")
}
