package main

import (
	"ShajeshJ/adventofcode_2022/day1"
	"log"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	sugar := logger.Sugar()
	day1.PartTwo(sugar)
}
