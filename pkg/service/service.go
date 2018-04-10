package service

import (
	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/logger"
	"github.com/takama/back-friend/pkg/logger/stdlog"
	"github.com/takama/back-friend/pkg/system"
	"github.com/takama/back-friend/pkg/version"
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

	// Wait signals
	signals := system.NewSignals()
	if err := signals.Wait(log, new(system.Handling)); err != nil {
		log.Error(err)
	}

	return nil
}
