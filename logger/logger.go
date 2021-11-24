package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Logger FtLogger
var instance string

func init() {
	godotenv.Load(".env")
	level, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		level = "debug"
	}
	instance, ok = os.LookupEnv("LOG_INSTANCE")
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

	Logger.logger = logger.WithFields(logrus.Fields{
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
	logger      *logrus.Entry
	logFile     *os.File
	logFilePath string

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

func (l *FtLogger) Close() {
	if l.isWriteToFile() {
		l.logFile.Close()

		if l.sdr != nil {
			err := l.sendLogs()
			l.logger.Errorf("send logs: %s", err)
		}
	}
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
		MachineID:  machineID(),
		LogContent: string(content),
		Instance:   instance,
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
	l.logger.Debug(args...)
}

func (l *FtLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *FtLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *FtLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *FtLogger) Warning(args ...interface{}) {
	l.logger.Warning(args...)
}

func (l *FtLogger) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(format, args...)
}

func (l *FtLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *FtLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *FtLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *FtLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *FtLogger) Print(args ...interface{}) {
	l.logger.Print(args...)
}

func (l *FtLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func machineID() string {
	machineID := ""
	if instance == "sender" || instance == "receiver" {
		machineID = "2"
	}

	return machineID
}
