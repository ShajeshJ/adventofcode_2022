package main

import (
	"log"

	"github.com/ShajeshJ/adventofcode_2022/solutions/day2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.ConsoleSeparator = " >> "
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}

	sugar := logger.Sugar()
	day2.PartTwo(sugar)
}
