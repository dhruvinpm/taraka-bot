package core

import (
	"context"
	"log"

	"github.com/dhruvinpm/taraka-bot/analyst"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type Pipeline struct {
	store   *memory.Store
	analyst *analyst.Analyst
	engine  *Engine
}

func NewPipeline(store *memory.Store, a *analyst.Analyst, engine *Engine) *Pipeline {
	return &Pipeline{store: store, analyst: a, engine: engine}
}

func (p *Pipeline) Run(ctx context.Context) error {
	log.Println("pipeline: starting")

	if _, err := p.engine.RunHunt(ctx, "", ""); err != nil {
		log.Printf("pipeline hunt error: %v", err)
	}

	if err := p.engine.RunAnalysis(ctx); err != nil {
		log.Printf("pipeline analysis error: %v", err)
	}

	if err := p.engine.RunOutreach(ctx); err != nil {
		log.Printf("pipeline outreach error: %v", err)
	}

	log.Println("pipeline: complete")
	return nil
}
