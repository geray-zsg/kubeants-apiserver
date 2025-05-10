package initiallize

import (
	"strings"

	"go.uber.org/zap/zapcore"
	"kubeants.io/config"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func InitLogger() {
	level := zapcore.InfoLevel
	switch strings.ToLower(config.CONF.Log.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "info":
		level = zapcore.InfoLevel
	}

	opts := zap.Options{
		Development: config.CONF.Log.Format == "console",
		TimeEncoder: zapcore.ISO8601TimeEncoder,
		Level:       level,
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
}
