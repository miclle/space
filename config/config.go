package config

import (
	"github.com/fox-gonic/fox/database"
	"github.com/fox-gonic/fox/logger"
)

// Configuration type
type Configuration struct {
	Addr     string           `mapstructure:"addr"`
	Secret   string           `mapstructure:"secret"`
	Env      string           `mapstructure:"env"`
	Logger   *logger.Config   `mapstructure:"logger"`
	Database *database.Config `mapstructure:"database"`
}
