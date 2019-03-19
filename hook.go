package hook

import (
	"fmt"
)

// List is just an abstraction over string array
type List []string

// Has checks if an item is in the list or not
func (l List) Has(s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}

// Actions is a list for valid actions
var Actions = List{
	"assigned",
	"unassigned",
	"review_requested",
	"review_request_removed",
	"opened",
	"edited",
	"closed",
	"reopened",
}

// Events is a list for valid events
var Events = List{
	"pull_request",
	"pull_request_review",
}

// Notifier interface needs to be implements by any struct so that it could be injected
// into Hook to send notifications
type Notifier interface {
	Notify(string, string) error
}

// Hook struct takes care of processing the payload and sends out notifications
type Hook struct {
	Notifier
	Payload *Payload
}

// Perform processes the payload and then notify based on the injected notifier
func (h *Hook) Perform() error {
	if !Actions.Has(h.Payload.Action) {
		return fmt.Errorf("unregistered action: %s", h.Payload.Action)
	}

	title, msg := h.Payload.Process()
	if title == "" || msg == "" {
		return fmt.Errorf("unable to process payload for action: %s", h.Payload.Action)
	}

	return h.Notify(title, msg)
}

// NewHook returns a new instance of Hook
func NewHook(p *Payload, n Notifier) *Hook {
	return &Hook{
		n,
		p,
	}
}
