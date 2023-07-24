package main

import (
	"flag"
	"os"
	"time"
	mytools "zapexample/tools"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int

const (
	Warn LogLevel = iota
	Info
	Debug
	Error
)
const dateTimeFormat = "[02/01/2006 15:04:05]"
const LOG_FILE = "log.txt"

var logger *zap.SugaredLogger = nil

func main() {

	interval := flag.Int("i", 5, "Greetings Interval")
	flag.Parse()

	initLogger()
	logger.Infow("Zap Package level logging example", "i", *interval)
	logger.Debug("Logger initialized")

	mytools.SetLogger(
		GetLogger(LogLevel(Debug)),
	)

	logger.Debug("Calling greetings from main")
	mytools.SayGreetings(*interval)

}

func initLogger() {
	// Initializing with default level for package main
	logger = GetLogger(LogLevel(Info))
}

func GetLogger(level LogLevel) *zap.SugaredLogger {
	// Configure logger options
	var l zapcore.Level
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(dateTimeFormat))
	}

	logFile, errLogFile := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if errLogFile != nil {
		defer logFile.Close()
		panic("Failed to open log file " + LOG_FILE)
	}

	// Set the initial logging level
	switch level {
	case Warn:
		l = zap.WarnLevel
	case Debug:
		l = zap.DebugLevel
	case Error:
		l = zap.ErrorLevel
	default:
		if logger == nil {
			l = zap.InfoLevel
		} else {
			return logger
		}
	}
	var err error
	_logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(logFile), l)).Sugar()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	// defer _logger.Sync()
	return _logger
}
