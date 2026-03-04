package browser

import (
	"fmt"
	"sync"

	"github.com/go-rod/rod"
)

type BrowserPool struct {
	mu       sync.Mutex
	browsers []*rod.Browser
	maxSize  int
	inUse    map[*rod.Browser]bool
}

func NewBrowserPool(maxSize int) *BrowserPool {
	return &BrowserPool{
		maxSize: maxSize,
		inUse:   make(map[*rod.Browser]bool),
	}
}

func (p *BrowserPool) Get() (*rod.Browser, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, b := range p.browsers {
		if !p.inUse[b] {
			p.inUse[b] = true
			return b, nil
		}
	}

	if len(p.browsers) < p.maxSize {
		b := rod.New().MustConnect()
		p.browsers = append(p.browsers, b)
		p.inUse[b] = true
		return b, nil
	}

	// Pool is at max capacity; return an error so callers can handle gracefully.
	return nil, fmt.Errorf("browser pool exhausted (max %d)", p.maxSize)
}

func (p *BrowserPool) Release(b *rod.Browser) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.inUse[b] = false
}

func (p *BrowserPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, b := range p.browsers {
		b.MustClose()
	}
	p.browsers = nil
}
