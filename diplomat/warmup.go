package diplomat

import (
	"time"
)

type WarmupManager struct {
	startPerDay      int
	incrementPerWeek int
	maxPerDay        int
	startDate        time.Time
}

func NewWarmupManager(startPerDay, incrementPerWeek, maxPerDay int) *WarmupManager {
	return &WarmupManager{
		startPerDay:      startPerDay,
		incrementPerWeek: incrementPerWeek,
		maxPerDay:        maxPerDay,
		startDate:        time.Now(),
	}
}

const warmupHoursPerWeek = 24 * 7

func (w *WarmupManager) GetDailyLimit() int {
	weeks := int(time.Since(w.startDate).Hours() / warmupHoursPerWeek)
	limit := w.startPerDay + (weeks * w.incrementPerWeek)
	if limit > w.maxPerDay {
		limit = w.maxPerDay
	}
	return limit
}
