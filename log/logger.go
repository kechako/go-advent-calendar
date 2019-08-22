package log

import (
	"log"

	stackdriver "github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap"
)

var (
	Logger = newLogger()
)

func newLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig = stackdriver.EncoderConfig
	l, err := config.Build()
	if err != nil {
		log.Fatal("failed to build new logger: ", err)
	}

	return l
}

