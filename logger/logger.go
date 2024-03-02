package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func ConstructLogger(cfg LoggerConfig) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	switch cfg.Environment {
	case Development:
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("new development logger")
		}
	case Production:
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("new production logger")
		}
	default:
		return nil, fmt.Errorf("unexpected environment for logger: %w", err)
	}

	defer logger.Sync()
	return logger, nil
}
