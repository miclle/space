package actions

import (
	"github.com/miclle/space/accounts"
	"github.com/miclle/space/config"
	"github.com/miclle/space/spaces"
)

// Actions type
type Actions struct {
	Configuration config.Configuration
	Accounter     accounts.Service
	Spacer        spaces.Service
}
