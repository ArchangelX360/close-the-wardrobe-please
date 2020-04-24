package main

import (
	"github.com/warthog618/gpiod"
	"time"
)

type Notifier struct {
	isNotifying bool
	cancel      *chan struct{}

	buzzer *gpiod.Line
}

func NewNotifier(c *gpiod.Chip, buzzerPin int) *Notifier {
	cancel := make(chan struct{}, 1)

	buzzer, err := c.RequestLine(buzzerPin, gpiod.AsOutput(1))
	if err != nil {
		errorLogger.Fatalf("failed to request line %d: %v", buzzerPin, err)
	}
	buzzer.SetValue(0)

	return &Notifier{
		isNotifying: false,
		cancel:      &cancel,
		buzzer:      buzzer,
	}
}

func (n *Notifier) fireNotificationIn(duration time.Duration) {
	select {
	case <-*n.cancel:
		n.buzzer.SetValue(0)
		n.isNotifying = false
		infoLogger.Println("notification cancelled")
	case <-time.After(duration):
		n.buzzer.SetValue(1)
		infoLogger.Println("notified!")
		n.isNotifying = false
	}
}

func (n *Notifier) flushCancellation() {
	for len(*n.cancel) > 0 {
		<-*n.cancel
	}
}

func (n *Notifier) Prepare() {
	infoLogger.Println("preparing to send notification...")
	n.flushCancellation()
	n.isNotifying = true
}

func (n *Notifier) FireCancellation() {
	if len(*n.cancel) == 0 {
		*n.cancel <- struct{}{}
	}
	infoLogger.Println("notification cancellation sent")
}
