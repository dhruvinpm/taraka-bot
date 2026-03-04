package core

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron   *cron.Cron
	engine *Engine
}

func NewScheduler(engine *Engine) *Scheduler {
	c := cron.New()
	return &Scheduler{cron: c, engine: engine}
}

func (s *Scheduler) Start() {
	ctx := context.Background()

	s.cron.AddFunc("0 * * * *", func() {
		if err := s.engine.RunPipeline(ctx); err != nil {
			log.Printf("scheduled pipeline error: %v", err)
		}
	})

	s.cron.AddFunc("0 8 * * *", func() {
		stats, err := s.engine.store.GetStats()
		if err != nil {
			log.Printf("stats error: %v", err)
			return
		}
		log.Printf("daily stats: %v", stats)
	})

	s.cron.Start()
	log.Println("scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
