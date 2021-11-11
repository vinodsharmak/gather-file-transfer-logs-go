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
	level, ok := os.LookupEnv("level")
	if !ok {
		level = "debug"
	}
	instance, ok := os.LookupEnv("instance")
	if !ok {
		instance = "anonymous raccoon"
	}

	logger := logrus.New()
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatalf("parse level: %s", err)
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
		logrus.Fatalf("prepare logfile: %s", err)
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

func (ft *FtLogger) Fatal(args ...interface{}) {
	ft.Logger.Fatal(args...)
}

func (ft *FtLogger) Fatalf(format string, args ...interface{}) {
	ft.Logger.Fatalf(format, args...)
}

func (ft *FtLogger) Print(args ...interface{}) {
	ft.Logger.Print(args...)
}

func (ft *FtLogger) Printf(format string, args ...interface{}) {
	ft.Logger.Printf(format, args...)
}
