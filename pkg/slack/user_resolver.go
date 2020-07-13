package slack

import (
	"context"
	goslack "github.com/slack-go/slack"
)

type UserResolver struct {
	client *goslack.Client
}

func NewUserResolver(token string) *UserResolver {
	return &UserResolver{
		client: goslack.New(token),
	}
}

// Notify resolves user by email, then sends him a message
func (n *UserResolver) ResolveUserEmail(ctx context.Context, userID string) (string, error) {
	user, err := n.client.GetUserInfoContext(ctx, userID)
	if err != nil {
		return "", err
	}

	return user.Profile.Email, nil
}
