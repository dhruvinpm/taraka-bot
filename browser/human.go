package browser

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

type HumanBehavior struct{}

func NewHumanBehavior() *HumanBehavior {
	return &HumanBehavior{}
}

func (h *HumanBehavior) RandomDelay(minMs, maxMs int) {
	d := minMs + rand.Intn(maxMs-minMs)
	time.Sleep(time.Duration(d) * time.Millisecond)
}

func (h *HumanBehavior) HumanType(el *rod.Element, text string) error {
	for _, ch := range text {
		if err := el.Type(input.Key(ch)); err != nil {
			return err
		}
		h.RandomDelay(50, 200)
	}
	return nil
}

func (h *HumanBehavior) HumanClick(el *rod.Element) error {
	h.RandomDelay(200, 800)
	return el.Click(proto.InputMouseButtonLeft, 1)
}
