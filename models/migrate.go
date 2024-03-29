package models

import (
	"github.com/fox-gonic/fox/database"
)

// Migrate models
func Migrate(db *database.Database) error {

	err := db.AutoMigrate(
		&Account{},
		&Authentication{},
		&Space{},
		&Page{},
		&PageContent{},
	)
	if err != nil {
		return err
	}

	return err
}
