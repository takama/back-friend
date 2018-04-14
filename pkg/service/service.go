package service

import (
	"fmt"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/back-friend/pkg/handlers"
	"github.com/takama/back-friend/pkg/logger"
	"github.com/takama/back-friend/pkg/logger/stdlog"
	"github.com/takama/back-friend/pkg/system"
	"github.com/takama/back-friend/pkg/version"

	"github.com/takama/bit"
)

// Run the service
func Run(cfg *config.Config) error {
	// Setup logger
	log := stdlog.New(&logger.Config{
		Level: cfg.LogLevel,
		Time:  true,
		UTC:   true,
	})

	log.Info("Version:", version.RELEASE)
	if cfg.LogLevel == logger.LevelDebug {
		log.Warnf("%s log level is used", logger.LevelDebug.String())
	}

	// Creates DB connection
	conn, name, err := db.New(cfg, log)
	if err != nil {
		log.Error("Could not establish database connection:", err)
		return err
	}
	log.Infof("Connection to database %s established successfully", name)

	// Define handlers
	h := handlers.New(conn)

	// Register new router
	r := bit.NewRouter()

	// Response for undefined methods
	r.SetupNotFoundHandler(h.Base(h.NotFound))

	// Configure router
	r.SetupMiddleware(h.Base)
	r.GET("/", h.Root)
	r.GET("/healthz", h.Health)
	r.GET("/readyz", h.Ready)
	r.GET("/info", h.Info)

	// Listen and serve handlers
	go r.Listen(fmt.Sprintf("%s:%d", cfg.LocalHost, cfg.LocalPort))

	// Wait signals
	signals := system.NewSignals()
	if err := signals.Wait(log, new(system.Handling)); err != nil {
		log.Error(err)
	}

	return nil
}
