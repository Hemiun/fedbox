package ecommerce

import (
	"git.sr.ht/~mariusor/lw"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/user"
	"github.com/go-ap/fedbox/internal/config"
)

var (
	userService *user.UserService
	logger      lw.Logger
	cfg         *config.Options
)

// New func init all required objects for ecommerce application
func New(ctl common.Control, db common.Storage, config *config.Options, l lw.Logger) {
	//ctl := cmd.New(db, *cfg, l)
	//_ = cmd.New(db, *cfg, l)

	logger = l
	cfg = config
	userService = user.NewUserService(ctl, db, cfg.BaseURL, logger)
}
