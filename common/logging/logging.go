package logging

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.ConsoleSeparator = " >> "
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger.Sugar()
}
