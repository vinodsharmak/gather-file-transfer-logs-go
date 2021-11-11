package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type ftLogger struct {
	Logger  *logrus.Entry
	logFile *os.File
}

var Logger *ftLogger

func init() {
	godotenv.Load(".env")
	level := os.Getenv("level")
	instance := os.Getenv("instance")

	logger := logrus.New()
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		Logger.Logger.Errorf("parse level: %s", err)
	}
	logger.SetLevel(lvl)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	Logger.Logger = logger.WithFields(logrus.Fields{
		"instance": instance,
	})

	logWriters := []io.Writer{os.Stdout}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		Logger.Errorf("Get user cache dir: %s", err)
	}
	err = os.Chdir(cacheDir)
	if err != nil {
		Logger.Errorf("Changing to chache directory: %s", err)
	}
	if _, err := os.Stat("cpaasFileTransfer"); os.IsNotExist(err) {
		err = os.Mkdir("cpaasFileTransfer", 0755)
		if err != nil {
			Logger.Errorf("Error in creating tmp directory: %s", err)
		}
	}

	logFilePath := filepath.Join(cacheDir, "cpaasFileTransfer", "fts.log")
	Logger.logFile, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		Logger.Logger.Fatal("Error in writing logs to logfile: %s", err)
	} else {
		logWriters = append(logWriters, Logger.logFile)
	}

	logger.SetOutput(io.MultiWriter(logWriters...))
}

func (ft *ftLogger) Close() {
	ft.logFile.Close()
}

func (ft *ftLogger) Debug(args ...interface{}) {
	ft.Logger.Debug(args...)
}

func (ft *ftLogger) Debugf(format string, args ...interface{}) {
	ft.Logger.Debugf(format, args...)
}

func (ft *ftLogger) Info(args ...interface{}) {
	ft.Logger.Info(args...)
}

func (ft *ftLogger) Infof(format string, args ...interface{}) {
	ft.Logger.Infof(format, args...)
}

func (ft *ftLogger) Warning(args ...interface{}) {
	ft.Logger.Warning(args...)
}

func (ft *ftLogger) Warningf(format string, args ...interface{}) {
	ft.Logger.Warningf(format, args...)
}

func (ft *ftLogger) Error(args ...interface{}) {
	ft.Logger.Error(args...)
}

func (ft *ftLogger) Errorf(format string, args ...interface{}) {
	ft.Logger.Errorf(format, args...)
}
