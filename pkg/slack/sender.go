package slack

import (
	"context"
	"github.com/kniepok/corgi"
	goslack "github.com/slack-go/slack"
)

type Sender struct {
	client *goslack.Client
}

func NewSender(token string) *Sender {
	return &Sender{client: goslack.New(token)}
}

// Notify resolves user by email, then sends him a message
func (n *Sender) Send(ctx context.Context, user corgi.User, message string) error {
	//channel, timestamp, text are ignored
	_, _, _, err := n.client.SendMessageContext(ctx, user.ID,
		goslack.MsgOptionText(message, false),
	)
	return err
}
