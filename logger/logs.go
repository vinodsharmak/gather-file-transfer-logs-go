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

	switch level {
	case "info":
		logrus.WithFields(logrus.Fields{
			"instance": instance,
		}).Info(s)

	case "error":
		logrus.WithFields(logrus.Fields{
			"instance": instance,
		}).Error(s)

	case "warning":
		logrus.WithFields(logrus.Fields{
			"instance": instance,
		}).Warning(s)

	case "debug":
		logrus.WithFields(logrus.Fields{
			"instance": instance,
		}).Debug(s)

	case "panic":
		logrus.WithFields(logrus.Fields{
			"instance": instance,
		}).Panic(s)

	default:
		logrus.WithFields(logrus.Fields{
			"instance": instance,
		}).Info(s)
	}

}
