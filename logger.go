package logger

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Entry

func init() {
	godotenv.Load(".env")
	level := os.Getenv("level")
	instance := os.Getenv("instance")

	logger := logrus.New()
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		Logger.Errorf("parse level: %s", err)
	}
	logger.SetLevel(lvl)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	Logger = logger.WithFields(logrus.Fields{
		"instance": instance,
	})
}
