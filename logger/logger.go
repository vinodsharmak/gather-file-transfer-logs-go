package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	Logger.instance, ok = os.LookupEnv("LOG_INSTANCE")
	if !ok {
		logrus.Warning("not found LOG_INSTANCE environment variable")
		Logger.instance = "anonymous raccoon"
	}

	Logger.logFilePath = os.Getenv("LOG_FILE_PATH")

	Logger.logger = logrus.New()
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatalf("parse level: %s", err)
	}
	Logger.logger.SetLevel(lvl)
	Logger.logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.logger.SetOutput(os.Stdout)

	Logger.loggerEntry = Logger.logger.WithFields(logrus.Fields{
		"instance": Logger.instance,
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

	Logger.logger.SetOutput(io.MultiWriter(logWriters...))
}

type FtLogger struct {
	logger      *logrus.Logger
	loggerEntry *logrus.Entry
	logFile     *os.File
	logFilePath string
	instance    string

	sdr *sender
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

func (l *FtLogger) SetInstance(instance string) {
	if instance == "" {
		logrus.Warning("instance: empty string")
		instance = "anonymous raccoon"
	}
	l.instance = instance
	l.loggerEntry = l.logger.WithFields(logrus.Fields{
		"instance": instance,
	})
}

func (l *FtLogger) SetLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatalf("parse level: %s", err)
	}
	l.logger.SetLevel(lvl)
}

func (l *FtLogger) SetLogFile(path string) {
	l.logFilePath = path
	logWriters := []io.Writer{os.Stdout}

	if l.isWriteToFile() {
		err := l.prepLogFile()
		if err != nil {
			logrus.Fatalf("prepare logfile: %s", err)
		} else {
			logWriters = append(logWriters, l.logFile)
		}
	}

	l.logger.SetOutput(io.MultiWriter(logWriters...))
}

func (l *FtLogger) Close() error {
	if l.isWriteToFile() {
		err := l.logFile.Close()
		if err != nil {
			msg := fmt.Sprintf("close log file: %s", err)
			logrus.Error(msg)

			return errors.New(msg)
		}

		if l.sdr != nil {
			err := l.sendLogs()
			if err != nil {
				msg := fmt.Sprintf("send logs: %s", err)
				logrus.Error(msg)

				return errors.New(msg)
			}
		}
	}

	return nil
}

func (l *FtLogger) sendLogs() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get user cache dir: %v", err)
	}
	logFilePath := filepath.Join(cacheDir, l.logFilePath)
	content, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		return fmt.Errorf("error in reading log file: %v", err)
	}

	reqData := request{
		MachineID:  l.machineID(),
		LogContent: string(content),
		Instance:   l.instance,
	}
	requestBody, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("request body: %v", err)
	}

	return l.sdr.send(requestBody)
}

func (l *FtLogger) SetSender(accessToken, url, machinePairID string) {
	l.sdr = newSender(accessToken, url, machinePairID)
}

func (l *FtLogger) Debug(args ...interface{}) {
	l.loggerEntry.Debug(args...)
}

func (l *FtLogger) Debugf(format string, args ...interface{}) {
	l.loggerEntry.Debugf(format, args...)
}

func (l *FtLogger) Info(args ...interface{}) {
	l.loggerEntry.Info(args...)
}

func (l *FtLogger) Infof(format string, args ...interface{}) {
	l.loggerEntry.Infof(format, args...)
}

func (l *FtLogger) Warning(args ...interface{}) {
	l.loggerEntry.Warning(args...)
}

func (l *FtLogger) Warningf(format string, args ...interface{}) {
	l.loggerEntry.Warningf(format, args...)
}

func (l *FtLogger) Error(args ...interface{}) {
	l.loggerEntry.Error(args...)
}

func (l *FtLogger) Errorf(format string, args ...interface{}) {
	l.loggerEntry.Errorf(format, args...)
}

func (l *FtLogger) Fatal(args ...interface{}) {
	l.loggerEntry.Fatal(args...)
}

func (l *FtLogger) Fatalf(format string, args ...interface{}) {
	l.loggerEntry.Fatalf(format, args...)
}

func (l *FtLogger) Print(args ...interface{}) {
	l.loggerEntry.Print(args...)
}

func (l *FtLogger) Printf(format string, args ...interface{}) {
	l.loggerEntry.Printf(format, args...)
}

func (l *FtLogger) machineID() string {
	machineID := ""
	if l.instance == "sender" || l.instance == "receiver" {
		machineID = "2"
	}

	return machineID
}
