package hook

import (
	"fmt"
	"os"
)

var (
	APIKey  string
	UserKey string
)

func init() {
	if os.Getenv("API_KEY") == "" || os.Getenv("USER_KEY") == "" {
		panic("keys are not set for API and recipient")
	}

	APIKey, UserKey = os.Getenv("API_KEY"), os.Getenv("USER_KEY")
}

type List []string

func (l List) Has(s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}

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

var Events = List{
	"pull_request",
	"pull_request_review",
}

type Notifier interface {
	Notify(string, string) error
}

type Hook struct {
	Notifier
	Payload *Payload
}

func (h *Hook) Perform() error {
	if !Actions.Has(h.Payload.Action) {
		return fmt.Errorf("unregistered action: %s", h.Payload.Action)
	}

	title, msg := h.Payload.Process()
	if title == "" || msg == "" {
		return fmt.Errorf("unable to process payload")
	}

	return h.Notify(title, msg)
}

func NewHook(p *Payload) *Hook {
	return &Hook{
		new(Pushover),
		p,
	}
}

func NewHookWithNotifier(p *Payload, n Notifier) *Hook {
	return &Hook{
		n,
		p,
	}
}
