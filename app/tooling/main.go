package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/halilylm/service/app/tooling/commands"
	"github.com/halilylm/service/business/sys/database"
	"go.uber.org/zap"
	"os"
)

var build = "develop"

type config struct {
	conf.Version
	Args conf.Args
	DB   struct {
		User         string `conf:"default:root"`
		Password     string `conf:"default:secret,mask"`
		Host         string `conf:"default:localhost"`
		Name         string `conf:"default:sales"`
		MaxIdleConns int    `conf:"default:2"`
		MaxOpenConns int    `conf:"default:0"`
		DisableTLS   bool   `conf:"default:true"`
	}
}

func main() {
	if err := run(zap.NewNop().Sugar()); err != nil {
		if !errors.Is(err, commands.ErrHelp) {
			fmt.Println("ERROR", err)
		}
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	cfg := config{
		Version: conf.Version{
			SVN:  build,
			Desc: "copyright information here",
		},
	}

	const prefix = "SALES"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		out, err := conf.String(&cfg)
		if err != nil {
			return fmt.Errorf("generating config for output: %w", err)
		}
		log.Infow("startup", "config", out)

		return fmt.Errorf("parsing config: %w", err)
	}

	return processCommands(cfg.Args, log, cfg)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(args conf.Args, log *zap.SugaredLogger, cfg config) error {
	dbConfig := database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	}

	switch args.Num(0) {
	case "migrate":
		if err := commands.Migrate(dbConfig); err != nil {
			return fmt.Errorf("migrating database: %w", err)
		}

	default:
		fmt.Println("migrate:    create the schema in the database")
		return commands.ErrHelp
	}

	return nil
}
