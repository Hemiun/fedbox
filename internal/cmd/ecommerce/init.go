package ecommerce

import (
	"git.sr.ht/~mariusor/lw"
	auth2 "github.com/go-ap/auth"
	"github.com/go-ap/client"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/mail"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/middleware"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/product"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/user"
	"github.com/go-ap/fedbox/internal/config"
)

var (
	userService    UserService
	productService ProductService
	mailService    MailService
	logger         lw.Logger
	cfg            *config.Options
)

// New func init all required objects for ecommerce application
func New(ctl common.Control, db common.Storage, config *config.Options, l lw.Logger) error {
	var err error
	logger = l
	cfg = config

	userService, err = user.NewUserService(ctl, db, cfg.BaseURL, logger)
	if err != nil {
		logger.Errorf("Can't init user service: %v", err)
		return err
	}

	productService = product.NewProductService(db, cfg.BaseURL, logger)

	cl := client.New(
		client.WithLogger(l.WithContext(lw.Ctx{"log": "client"})),
		client.SkipTLSValidation(!cfg.Env.IsProd()),
	)
	middleware.AuthService, err = auth2.New(auth2.WithURL(cfg.BaseURL),
		auth2.WithStorage(db),
		auth2.WithClient(cl), //TODO:
		auth2.WithLogger(l),
	)
	if err != nil {
		logger.Errorf("Can't init auth service: %v", err)
		return err
	}

	//email service initialization
	mailService = mail.NewMailer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUser, cfg.SmtpPass, cfg.SmtpFrom, logger)

	return nil
}
