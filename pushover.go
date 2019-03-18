package hook

import "github.com/gregdel/pushover"

type Pushover struct{}

func (p *Pushover) Notify(title, msg string) error {
	app := pushover.New(APIKey)
	recipient := pushover.NewRecipient(UserKey)
	message := pushover.NewMessageWithTitle(msg, title)
	_, err := app.SendMessage(message, recipient)
	return err
}
