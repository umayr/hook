package hook

import (
	"github.com/gregdel/pushover"
)

// Pushover is an struct that implements Notifier interface
type Pushover struct {
	APIKey  string
	UserKey string
}

// Notify sends out the notification using Pushover API.
// It results in error in case there's missing or bad API or User key.
func (p *Pushover) Notify(title, msg string) error {
	app := pushover.New(p.APIKey)
	recipient := pushover.NewRecipient(p.UserKey)
	message := pushover.NewMessageWithTitle(msg, title)
	_, err := app.SendMessage(message, recipient)
	return err
}

// NewPushover returns a new instance of Pushover struct
func NewPushover(APIKey, UserKey string) *Pushover {
	return &Pushover{APIKey, UserKey}
}
