package logger

import (
	"os"

	"bitbucket.org/gath3rio/gather-service-logs/constants"
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

	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		err = os.Mkdir("tmp", 0755)
		if err != nil {
			Logger.Errorf("Error in creating tmp directory: %s", err)
		}
	}
	file, err := os.OpenFile(constants.LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal("Error in writing logs to file: %s", err)
	}
	logger.SetOutput(file)
}
