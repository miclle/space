package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/fox-gonic/fox/configurations"
	"github.com/fox-gonic/fox/database"
	"github.com/fox-gonic/fox/engine"
	"github.com/fox-gonic/fox/logger"
	"golang.org/x/sync/errgroup"

	"github.com/miclle/space/accounts"
	"github.com/miclle/space/config"
	"github.com/miclle/space/models"
	"github.com/miclle/space/spaces"
)

var (
	g errgroup.Group
)

func main() {

	log := logger.New("Initialize")

	var confPath string
	flag.StringVar(&confPath, "f", "", "-f=/path/to/config")
	flag.Parse()

	var configuration config.Configuration
	err := configurations.Parse(confPath, &configuration)
	if err != nil {
		log.Fatalf("parse config failed, err: %+v", err)
	}

	database, err := database.New(configuration.Database)
	if err != nil {
		log.Fatalf("database init failed, err: %+v", err)
	}

	err = models.Migrate(database)
	if err != nil {
		log.Fatalf("migrate models failed, err: %+v", err)
	}

	accounter, err := accounts.NewService(database)
	if err != nil {
		log.Fatalf("new accounts service failed, err: %+v", err)
	}

	spacer, err := spaces.NewService(database)
	if err != nil {
		log.Fatalf("new spaces service failed, err: %+v", err)
	}

	engine.SetMode(configuration.Env)

	platformServer := &http.Server{
		Addr:         configuration.Addr,
		Handler:      router(configuration, accounter, spacer),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return platformServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
