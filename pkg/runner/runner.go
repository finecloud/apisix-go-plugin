package runner

import (
	"os"

	"go.uber.org/zap/zapcore"

	"github.com/finecloud/apisix-oauth2-plugin/internal/server"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/log"
)

// RunnerConfig is the configuration of the runner
type RunnerConfig struct {
	// LogLevel is the level of log, default to `zapcore.InfoLevel`
	LogLevel zapcore.Level
	// LogOutput is the output of log, default to `os.Stdout`
	LogOutput zapcore.WriteSyncer
}

// Run starts the runner and listen the socket configured by environment variable "APISIX_LISTEN_ADDRESS"
func Run(cfg RunnerConfig) {
	if cfg.LogOutput == nil {
		cfg.LogOutput = os.Stdout
	}
	log.NewLogger(cfg.LogLevel, cfg.LogOutput)
	server.Run()
}
