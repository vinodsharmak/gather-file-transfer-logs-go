package logger

import (
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Logger FtLogger

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

	Logger.logFile, err = prepLogFile("cpaasFileTransfer", "fts.log")
	if err != nil {
		Logger.Logger.Fatal("Error in writing logs to logfile: %s", err)
	} else {
		logWriters = append(logWriters, Logger.logFile)
	}

	logger.SetOutput(io.MultiWriter(logWriters...))
}

type FtLogger struct {
	Logger  *logrus.Entry
	logFile *os.File
}

func (ft *FtLogger) Close() {
	ft.logFile.Close()
}

func (ft *FtLogger) Debug(args ...interface{}) {
	ft.Logger.Debug(args...)
}

func (ft *FtLogger) Debugf(format string, args ...interface{}) {
	ft.Logger.Debugf(format, args...)
}

func (ft *FtLogger) Info(args ...interface{}) {
	ft.Logger.Info(args...)
}

func (ft *FtLogger) Infof(format string, args ...interface{}) {
	ft.Logger.Infof(format, args...)
}

func (ft *FtLogger) Warning(args ...interface{}) {
	ft.Logger.Warning(args...)
}

func (ft *FtLogger) Warningf(format string, args ...interface{}) {
	ft.Logger.Warningf(format, args...)
}

func (ft *FtLogger) Error(args ...interface{}) {
	ft.Logger.Error(args...)
}

func (ft *FtLogger) Errorf(format string, args ...interface{}) {
	ft.Logger.Errorf(format, args...)
}