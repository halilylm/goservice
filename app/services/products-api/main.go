package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/halilylm/service/app/services/products-api/handlers"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var build = "develop"

func main() {
	log, err := initLogger("service-api")
	if err != nil {
		fmt.Println("error constructing logger")
		os.Exit(1)
	}
	defer log.Sync()
	log.Info("hello")

	// =========================================
	// Configuration
	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "copyright of this app",
		},
	}

	const prefix = "PRODUCT"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	debugApp := fiber.New()
	handlers.DebugRoutes(debugApp, build, log)

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	})

	// Default middleware config
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(pprof.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello there")
	})

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- debugApp.Listen(cfg.Web.DebugHost)
		serverErrors <- app.Listen(cfg.Web.APIHost)
	}()

	select {
	case err := <-serverErrors:
		log.Fatalw("server error", "error", err.Error())
	case sig := <-quit:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown completed", "signal", sig)

		if err := app.ShutdownWithTimeout(cfg.Web.ShutdownTimeout); err != nil {
			log.Fatalw("error shutting down the server gracefully", "err", err.Error())
		}
	}
}

func initLogger(service string) (*zap.SugaredLogger, error) {
	// Construct the application logger.
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"service": "products-api",
	}
	log, err := config.Build()
	if err != nil {
		return nil, err
	}
	return log.Sugar(), nil
}
