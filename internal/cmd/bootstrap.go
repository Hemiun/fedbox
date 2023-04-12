package cmd

import (
	"fmt"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/bootstrap"
	"os"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/errors"
	"github.com/go-ap/fedbox"
	ap "github.com/go-ap/fedbox/activitypub"
	"github.com/go-ap/fedbox/internal/config"
	s "github.com/go-ap/fedbox/storage"
	"github.com/urfave/cli/v2"
)

var BootstrapCmd = &cli.Command{
	Name:  "bootstrap",
	Usage: "Bootstrap a new postgres or bolt database helper",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "root",
			Usage: "root account of postgres server (default: postgres)",
			Value: "postgres",
		},
		&cli.StringFlag{
			Name:  "sql",
			Usage: "path to the queries for initializing the database",
			Value: "postgres",
		},
		&cli.StringFlag{
			Name:  "key-type",
			Usage: fmt.Sprintf("Type of keys to generate: %v", []string{fedbox.KeyTypeED25519, fedbox.KeyTypeRSA}),
			Value: fedbox.KeyTypeED25519,
		},
	},
	Action: bootstrapAct(&ctl),
	Subcommands: []*cli.Command{
		reset,
		initApp,
	},
}

var reset = &cli.Command{
	Name:   "reset",
	Usage:  "reset an existing database",
	Action: resetAct(&ctl),
}

func resetAct(c *Control) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		err := Reset(c.Conf)
		if err != nil {
			return err
		}
		return Bootstrap(c.Conf, c.Service)
	}
}

func bootstrapAct(c *Control) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		keyType := ctx.String("keyType")
		if err := Bootstrap(c.Conf, c.Service); err != nil {
			Errf("Error adding service: %s\n", err)
			return err
		}
		c.Service = ap.Self(ap.DefaultServiceIRI(c.Conf.BaseURL))
		if _, err := c.Storage.Save(&c.Service); err != nil {
			Errf("Error adding service: %s\n", err)
			return err
		}
		if metaSaver, ok := ctl.Storage.(s.MetadataTyper); ok {
			if err := AddKeyToItem(metaSaver, &c.Service, keyType); err != nil {
				Errf("Error saving metadata for service: %s", err)
				return err
			}
		}
		return nil
	}
}

type storageConf struct {
	Storage     config.StorageType
	Path        string
	BaseURL     string
	CacheEnable bool
}

func confFn(opt config.Options) storageConf {
	return storageConf{Path: opt.BaseStoragePath(), CacheEnable: opt.StorageCache, BaseURL: opt.BaseURL, Storage: opt.Storage}
}

func Bootstrap(conf config.Options, service vocab.Item) error {
	if err := bootstrapFn(confFn(conf), service); err != nil {
		return errors.Annotatef(err, "Unable to create %s db for storage %s", conf.BaseStoragePath(), conf.Storage)
	}
	fmt.Fprintf(os.Stdout, "Successfuly created %s db for storage %s\n", conf.BaseStoragePath(), conf.Storage)
	return nil
}

func Reset(conf config.Options) error {
	if err := cleanFn(confFn(conf)); err != nil {
		return errors.Annotatef(err, "Unable to reset %s db for storage %s", conf.BaseStoragePath(), conf.Storage)
	}
	fmt.Fprintf(os.Stdout, "Successful reset %s db for storage %s\n", conf.BaseStoragePath(), conf.Storage)
	return nil
}

var initApp = &cli.Command{
	Name:  "initApp",
	Usage: "app initialization: new oauth client, new superuser",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "redirectUri",
			Value:    nil,
			Usage:    "The redirect URIs for current application",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "secret",
			Value:    "",
			Usage:    "The client secret",
			Required: true,
		},
	},
	Action: initAppAct(&ctl),
}

func initAppAct(c *Control) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return bootstrap.InitNewApp(ctx, &ctl, c.Conf.BaseURL, ctl.Storage)
	}
}
