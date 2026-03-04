package core

import (
	"context"
	"fmt"
	"log"

	"github.com/dhruvinpm/taraka-bot/analyst"
	"github.com/dhruvinpm/taraka-bot/browser"
	"github.com/dhruvinpm/taraka-bot/diplomat"
	localmail "github.com/dhruvinpm/taraka-bot/gmail"
	"github.com/dhruvinpm/taraka-bot/hunter"
	"github.com/dhruvinpm/taraka-bot/llm"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type Config struct {
	DBPath           string
	GeminiAPIKey     string
	GeminiModel      string
	GeminiDailyLimit int
	LocalLLMURL      string
	LocalLLMModel    string
	SMTPConfig       localmail.SMTPConfig
	GmailAddress     string
	GmailCredentials string
	GmailToken       string
	WarmupStart      int
	WarmupInc        int
	WarmupMax        int
	VaultPath        string
	VaultKey         string
	DefaultCountry   string
	DefaultNiche     string
}

type Engine struct {
	cfg       Config
	store     *memory.Store
	llmRouter *llm.Router
	hunter    *hunter.Hunter
	analyst   *analyst.Analyst
	diplomat  *diplomat.Diplomat
	pool      *browser.BrowserPool
	pipeline  *Pipeline
	scheduler *Scheduler
	paused    bool
}

func New(ctx context.Context, cfg Config) (*Engine, error) {
	store, err := memory.NewStore(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	var primaryLLM llm.Provider
	if cfg.GeminiAPIKey != "" {
		g, err := llm.NewGeminiProvider(ctx, cfg.GeminiAPIKey, cfg.GeminiModel, cfg.GeminiDailyLimit)
		if err != nil {
			log.Printf("gemini init error: %v", err)
		} else {
			primaryLLM = g
		}
	}

	var fallbackLLM llm.Provider
	if cfg.LocalLLMURL != "" {
		fallbackLLM = llm.NewLocalProvider(cfg.LocalLLMURL, cfg.LocalLLMModel)
	}

	router := llm.NewRouter(primaryLLM, fallbackLLM)

	pool := browser.NewBrowserPool(3)

	h := hunter.NewHunter(store, pool, primaryLLM)
	a := analyst.NewAnalyst(store)

	var gmailClient *localmail.Client
	if cfg.GmailCredentials != "" && cfg.GmailToken != "" {
		gmailClient, err = localmail.NewClient(ctx, cfg.GmailCredentials, cfg.GmailToken, cfg.GmailAddress)
		if err != nil {
			log.Printf("gmail client error: %v", err)
		}
	}

	if gmailClient == nil {
		gmailClient = localmail.NewNullClient()
	}

	dip := diplomat.NewDiplomat(store, gmailClient, diplomat.Config{
		SMTPConfig:  cfg.SMTPConfig,
		WarmupStart: cfg.WarmupStart,
		WarmupInc:   cfg.WarmupInc,
		WarmupMax:   cfg.WarmupMax,
	}, router)

	e := &Engine{
		cfg:       cfg,
		store:     store,
		llmRouter: router,
		hunter:    h,
		analyst:   a,
		diplomat:  dip,
		pool:      pool,
	}

	e.pipeline = NewPipeline(store, a, e)
	e.scheduler = NewScheduler(e)

	return e, nil
}

func (e *Engine) Start(ctx context.Context) {
	e.scheduler.Start()
	log.Println("engine started")
}

func (e *Engine) Stop() {
	e.scheduler.Stop()
	e.pool.Close()
	e.store.Close()
	log.Println("engine stopped")
}

func (e *Engine) Pause() {
	e.paused = true
}

func (e *Engine) Resume() {
	e.paused = false
}

func (e *Engine) IsPaused() bool {
	return e.paused
}

// RunHunt runs the full hunt → analyse → outreach pipeline, sending progress
// updates via the optional notify callback (used by the Telegram handler).
func (e *Engine) RunHunt(ctx context.Context, country, niche string, notify func(string)) error {
	if e.paused {
		return nil
	}
	if country == "" {
		country = e.cfg.DefaultCountry
	}
	if niche == "" {
		niche = e.cfg.DefaultNiche
	}

	send := func(msg string) {
		log.Println(msg)
		if notify != nil {
			notify(msg)
		}
	}

	// --- Step 1: Hunt ---
	send("🔍 Hunting leads...")
	cfg := hunter.HuntConfig{
		Country:  country,
		Niche:    niche,
		Mode:     hunter.Gentle,
		MaxLeads: 50,
	}
	result, err := e.hunter.Hunt(ctx, cfg)
	if err != nil {
		return err
	}

	for _, srcErr := range result.Errors {
		send(fmt.Sprintf("⚠️ Source error: %s", srcErr))
	}

	rawCount := len(result.Leads)
	if rawCount == 0 {
		send("❌ No leads found. Try different keywords.")
		return nil
	}
	send(fmt.Sprintf("📋 Found %d raw leads", rawCount))

	// --- Step 2: Analyse / verify ---
	send("🔬 Verifying and scoring leads...")
	leads, err := e.store.GetLeadsByStatus("raw", cfg.MaxLeads)
	if err != nil {
		return err
	}
	verifiedCount := 0
	for _, lead := range leads {
		enriched := e.analyst.Analyze(lead)
		if err := e.store.UpsertLead(enriched); err != nil {
			log.Printf("upsert error: %v", err)
		}
		if enriched.Status == "verified" {
			verifiedCount++
		}
	}
	send(fmt.Sprintf("✅ %d verified leads (with valid emails)", verifiedCount))

	if verifiedCount == 0 {
		send("ℹ️ No leads have verified email addresses yet — check back after sourcing more leads.")
		return nil
	}

	// --- Step 3: Outreach ---
	send("✍️ Composing and sending outreach emails...")
	verifiedLeads, err := e.store.GetLeadsByStatus("verified", cfg.MaxLeads)
	if err != nil {
		return err
	}
	sentCount := 0
	for _, lead := range verifiedLeads {
		if err := e.diplomat.SendInitialOutreach(ctx, lead); err != nil {
			log.Printf("outreach error for %s: %v", lead.Email, err)
		} else {
			sentCount++
		}
	}
	send(fmt.Sprintf("📧 Sent %d outreach emails", sentCount))

	if err := e.diplomat.ProcessFollowUps(ctx); err != nil {
		log.Printf("follow-up error: %v", err)
	}

	return nil
}

func (e *Engine) RunAnalysis(ctx context.Context) error {
	if e.paused {
		return nil
	}
	leads, err := e.store.GetLeadsByStatus("raw", 100)
	if err != nil {
		return err
	}
	for _, lead := range leads {
		enriched := e.analyst.Analyze(lead)
		if err := e.store.UpsertLead(enriched); err != nil {
			log.Printf("upsert error: %v", err)
		}
	}
	return nil
}

func (e *Engine) RunOutreach(ctx context.Context) error {
	if e.paused {
		return nil
	}
	leads, err := e.store.GetLeadsByStatus("verified", 10)
	if err != nil {
		return err
	}
	for _, lead := range leads {
		if err := e.diplomat.SendInitialOutreach(ctx, lead); err != nil {
			log.Printf("outreach error for %s: %v", lead.Email, err)
		}
	}
	return e.diplomat.ProcessFollowUps(ctx)
}

func (e *Engine) RunPipeline(ctx context.Context) error {
	return e.pipeline.Run(ctx)
}

func (e *Engine) GetStore() *memory.Store {
	return e.store
}

func (e *Engine) GetLLM() *llm.Router {
	return e.llmRouter
}
