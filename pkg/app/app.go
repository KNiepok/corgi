package app

import (
	"context"
	"github.com/go-chi/chi"
	cli "github.com/jawher/mow.cli"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

// Application wraps common logic for booting application
type Application struct {
	Cli    *cli.Cli
	Setup  func(ctx context.Context)
	Config interface{}
	Router *chi.Mux
}

// NewApplication creates new application with given config
func NewApplication(name string, description string, conf interface{}) *Application {
	app := &Application{
		Cli:    cli.App(name, description),
		Setup:  func(ctx context.Context) {},
		Config: conf,
		Router: chi.NewRouter(),
	}

	app.Cli.Action = func() {
		listenAddr := app.Cli.String(cli.StringOpt{
			Name:   "listen-http",
			Desc:   "Listen address for HTTP server",
			Value:  ":8080",
			EnvVar: "LISTEN",
		})
		if err := envconfig.Process("", app.Config); err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		app.Setup(ctx)
		log.Fatal(http.ListenAndServe(*listenAddr, app.Router))
	}

	return app
}

// Run boots application
func (app *Application) Run(params []string) error {
	return app.Cli.Run(params)
}
