package logger

import (
	"io"
	"os"
	"path/filepath"

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

	logWriters := []io.Writer{os.Stdout}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		logger.Errorf("Get user cache dir: %s", err)
	}
	logFilePath := filepath.Join(cacheDir, "cpaasFileTransfer", "fts.log")
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal("Error in writing logs to logfile: %s", err)
	} else {
		logWriters = append(logWriters, file)
	}
	logger.SetOutput(io.MultiWriter(logWriters...))
}
