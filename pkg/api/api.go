package api

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/kniepok/corgi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"net/http"
)

type API struct {
	verificationToken string
	service           corgi.SubscriptionService
	userResolver      corgi.UserResolver
	logger            *logrus.Logger
}

func NewAPI(
	service corgi.SubscriptionService,
	userResolver corgi.UserResolver,
	verificationToken string,
) *API {
	return &API{
		service:           service,
		userResolver:      userResolver,
		verificationToken: verificationToken,
		logger:            logrus.New(),
	}
}

func (api *API) Mount(router *chi.Mux) {
	router.Post("/subscribe", api.handleSubscribe)
	router.Post("/unsubscribe", api.handleUnsubscribe)
}

func (api *API) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !s.ValidateToken(api.verificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if s.Command != "/subscribe" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := api.resolveUser(ctx, s.UserID)
	if err != nil {
		api.logger.Errorf("failed to resolve user with id = %s", s.UserID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subDetails, err := corgi.NewSubscriptionDetails(s.Text)
	if err != nil {
		api.logger.Infof("failed to parse interval = %s", s.Text)
		writeResponse(w,
			errorIntro,
			subscriptionErrorExplanation(s.Text),
			example,
		)
		return
	}

	err = api.service.Subscribe(ctx, corgi.Subscription{
		User:    user,
		Details: subDetails,
	})
	if err != nil {
		switch errors.Cause(err) {
		case corgi.ErrInvalidInterval:
			// todo remove that or what
			api.logger.Infof("failed to parse interval = %s", s.Text)
			writeResponse(w,
				errorIntro,
				subscriptionErrorExplanation(s.Text),
				example,
			)
			return
		}
		api.logger.Errorf("subscribe api: unexpected err = %s", err.Error())
		writeResponse(w, unexpectedError)
		return
	}
	writeResponse(w, subscribed)
}

func (api *API) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !s.ValidateToken(api.verificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if s.Command != "/unsubscribe" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := api.resolveUser(ctx, s.UserID)
	if err != nil {
		api.writeError(w, err)
		return
	}

	if err := api.service.Unsubscribe(ctx, user); err != nil {
		api.writeError(w, err)
		return
	}
	writeResponse(w, unsubscribed)
}

func (api *API) resolveUser(ctx context.Context, userID string) (corgi.User, error) {
	userEmail, err := api.userResolver.ResolveUserEmail(ctx, userID)
	if err != nil {
		return corgi.User{}, err
	}
	return corgi.User{
		ID:    userID,
		Email: userEmail,
	}, nil
}

func writeResponse(w http.ResponseWriter, blocks ...block) {
	message, _ := json.Marshal(slackResponse{
		Blocks: blocks,
	})
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(message)
}

func (api *API) writeError(w http.ResponseWriter, err error) {
	api.logger.Errorf("unsubscribe api: unexpected err = %s", err.Error())
	writeResponse(w, unexpectedError)
}
