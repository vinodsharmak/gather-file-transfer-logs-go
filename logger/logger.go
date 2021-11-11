package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Logger FtLogger

func init() {
	godotenv.Load(".env")
	level, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		level = "debug"
	}
	instance, ok := os.LookupEnv("LOG_INSTANCE")
	if !ok {
		logrus.Warning("not found LOG_INSTANCE environment variable")
		instance = "anonymous raccoon"
	}

	Logger.logFilePath = os.Getenv("LOG_FILE_PATH")

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

	if Logger.isWriteToFile() {
		err = Logger.prepLogFile()
		if err != nil {
			logrus.Fatalf("prepare logfile: %s", err)
		} else {
			logWriters = append(logWriters, Logger.logFile)
		}
	}

	logger.SetOutput(io.MultiWriter(logWriters...))
}

type FtLogger struct {
	Logger      *logrus.Entry
	logFile     *os.File
	logFilePath string
}

func (l *FtLogger) isWriteToFile() bool {
	return len(l.logFilePath) > 0
}

func (l *FtLogger) prepLogFile() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get user cache dir: %s", err)
	}

	logFileDir := filepath.Dir(l.logFilePath)
	logFileDirPath := filepath.Join(cacheDir, logFileDir)
	err = os.MkdirAll(logFileDirPath, 0755)
	if err != nil {
		return fmt.Errorf("make \"%s\" dir: %s", logFileDirPath, err)
	}

	logFilePath := filepath.Join(cacheDir, l.logFilePath)

	l.logFile, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (l *FtLogger) Close() {
	if l.isWriteToFile() {
		l.logFile.Close()
	}
}

func (l *FtLogger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *FtLogger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

func (l *FtLogger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *FtLogger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *FtLogger) Warning(args ...interface{}) {
	l.Logger.Warning(args...)
}

func (l *FtLogger) Warningf(format string, args ...interface{}) {
	l.Logger.Warningf(format, args...)
}

func (l *FtLogger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *FtLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func (l *FtLogger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

func (l *FtLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

func (l *FtLogger) Print(args ...interface{}) {
	l.Logger.Print(args...)
}

func (l *FtLogger) Printf(format string, args ...interface{}) {
	l.Logger.Printf(format, args...)
}
