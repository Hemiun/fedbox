package ecommerce

import (
	"git.sr.ht/~mariusor/lw"
	"github.com/go-ap/fedbox/ecommerce/user"
	"github.com/go-ap/fedbox/internal/cmd"
	"github.com/go-ap/fedbox/internal/config"
)

var (
	userService *user.UserService
	logger      lw.Logger
	cfg         *config.Options
)

func New(ctl *cmd.Control, config *config.Options, l lw.Logger) {
	logger = l
	cfg = config
	userService = user.NewUserService(ctl, cfg.BaseURL, logger)
}
