package logger

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}
func Log(s string) {
	godotenv.Load(".env")
	level := os.Getenv("level")
	instance := os.Getenv("instance")

	contextLogger := logrus.WithFields(logrus.Fields{
		"instance": instance,
	})

	switch level {
	case "info":
		contextLogger.Info(s)

	case "error":
		contextLogger.Error(s)

	case "warning":
		contextLogger.Warning(s)

	case "debug":
		contextLogger.Debug(s)

	case "panic":
		contextLogger.Panic(s)

	default:
		contextLogger.Info(s)
	}

}
