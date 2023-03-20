package main

import (
	"flag"

	"github.com/fox-gonic/fox/configurations"
	"github.com/fox-gonic/fox/database"
	"github.com/fox-gonic/fox/engine"
	"github.com/fox-gonic/fox/logger"

	"github.com/miclle/space/config"
	"github.com/miclle/space/models"
)

func main() {

	log := logger.New("Initialize")

	var confPath string
	flag.StringVar(&confPath, "f", "", "-f=/path/to/config")
	flag.Parse()

	var config config.Configuration
	err := configurations.Parse(confPath, &config)
	if err != nil {
		log.Fatalf("parse config failed, err: %+v", err)
	}

	database, err := database.New(config.Database)
	if err != nil {
		log.Fatalf("database init failed, err: %+v", err)
	}

	err = models.Migrate(database)
	if err != nil {
		log.Fatalf("migrate models failed, err: %+v", err)
	}

	engine.SetMode(config.Env)

}
