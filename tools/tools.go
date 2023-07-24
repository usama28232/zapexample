package tools

import (
	"time"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func SayGreetings(i int) {

	logger.Debug("Got a greeting request")

	// logger.Warn("proceeding with cold heart greetings: Hi")
	// logger.Error("Cold Hearted Greetings, Hi")

	logger.Debugw("Processing Greeting with interval", "i", i)

	ticker := time.NewTicker(time.Second * time.Duration(i))

	for range ticker.C {
		logger.Info("Greetings brother")
	}
}

func SetLogger(l *zap.SugaredLogger) {
	logger = l
}
