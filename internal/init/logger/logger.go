package logger

import (
	"os"

	"github.com/mojganchakeri/whatsapp-manager/internal/init/config"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
	Panic(msg string)
}

type logrusLogger struct {
	logger *logrus.Logger
}

func New(cfg config.Config) Logger {

	customFormatter := logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}

	var logger = logrus.New()

	logger.SetFormatter(&customFormatter)
	logger.Out = os.Stdout

	switch cfg.GetConfig().LogLevel {
	case "info":
		logger.Level = logrus.InfoLevel
	case "warn":
		logger.Level = logrus.WarnLevel
	case "error":
		logger.Level = logrus.ErrorLevel
	case "fatal":
		logger.Level = logrus.FatalLevel
	case "panic":
		logger.Level = logrus.PanicLevel
	case "debug":
		logger.Level = logrus.DebugLevel
	case "trace":
		logger.Level = logrus.TraceLevel
	default:
		logger.Level = logrus.InfoLevel
	}

	return &logrusLogger{
		logger: logger,
	}
}

func (l logrusLogger) Info(msg string) {
	l.logger.Infoln(msg)
}

func (l logrusLogger) Warn(msg string) {
	l.logger.Warnln(msg)
}

func (l logrusLogger) Error(msg string) {
	l.logger.Errorln(msg)
}

func (l logrusLogger) Debug(msg string) {
	l.logger.Debugln(msg)
}

func (l logrusLogger) Panic(msg string) {
	l.logger.Panicln(msg)
}
