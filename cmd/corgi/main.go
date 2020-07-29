package main

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kniepok/corgi/pkg/api"
	"github.com/kniepok/corgi/pkg/app"
	"github.com/kniepok/corgi/pkg/cron"
	"github.com/kniepok/corgi/pkg/service"
	"github.com/kniepok/corgi/pkg/slack"
	"github.com/kniepok/corgi/pkg/sqlite"
	"github.com/kniepok/corgi/pkg/tempo"
	"os"
)

type config struct {
	SlackToken             string `envconfig:"SLACK_TOKEN" default:""`
	SlackVerificationToken string `envconfig:"SLACK_VERIFICATION_TOKEN" default:""`
	DebugMode              bool   `envconfig:"DEBUG_MODE" default:"false"`
	TempoToken             string `envconfig:"TEMPO_TOKEN" default:""`
	JiraToken              string `envconfig:"JIRA_TOKEN" default:""`
	JiraUsername           string `envconfig:"JIRA_USERNAME" default:""`
	JiraBaseURL            string `envconfig:"JIRA_BASE_URL" default:"https://gogoapps.atlassian.net"`
	TempoBaseURL           string `envconfig:"TEMPO_BASE_URL" default:"https://api.tempo.io/core/3"`
}

func main() {
	var conf config
	application := app.NewApplication(os.Args[0], "API for user authentication", &conf)
	application.Setup = func(ctx context.Context) {
		storage := getStorage(conf)
		msgGenerator, err := tempo.NewNotificationGenerator(
			conf.TempoToken,
			conf.JiraToken,
			conf.JiraUsername,
			conf.JiraBaseURL,
			conf.TempoBaseURL,
		)
		if err != nil {
			panic(err)
		}
		notificationService := service.NewNotificationService(
			msgGenerator,
			slack.NewSender(conf.SlackToken),
		)
		subscriptionService := service.NewSubscriptionService(
			storage,
			cron.NewScheduler(),
			notificationService,
		)

		userResolver := slack.NewUserResolver(conf.SlackToken)
		api.NewAPI(
			subscriptionService,
			userResolver,
			conf.SlackVerificationToken,
		).Mount(application.Router)

		starter := service.NewStarterService(
			subscriptionService, storage)
		err = starter.Start(ctx)
		if err != nil {
			panic(err)
		}
	}

	err := application.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func getStorage(conf config) *sqlite.SubscriptionStorage {
	db, err := gorm.Open("sqlite3", "../../tmp/gorm.db")
	if err != nil {
		panic(err)
	}
	if conf.DebugMode {
		db = db.Debug()
	}
	return sqlite.NewSubscriptionStorage(db)
}
