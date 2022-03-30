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
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Warnln("unable to read .env, falling back to os variables")
	}
	Logger.logger = logrus.New()

	err = Logger.SetLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logrus.Errorf("set level: %s", err)
	}

	Logger.logger.SetFormatter(&logrus.JSONFormatter{})

	Logger.SetInstance(os.Getenv("LOG_INSTANCE"))

	err = Logger.SetLoggerOutput(os.Getenv("LOG_FILE_PATH"))
	if err != nil {
		logrus.Errorf("set log file: %s", err)
	}
}

type FtLogger struct {
	logger      *logrus.Logger
	loggerEntry *logrus.Entry
	logFile     *os.File
	logFilePath string
	instance    string
	debugger    bool
	sdr         *sender
}

func (l *FtLogger) isWriteToFile() bool {
	return len(l.logFilePath) > 0
}

func (l *FtLogger) prepLogFile() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get user cache dir: %w", err)
	}

	logFileDir := filepath.Dir(l.logFilePath)
	logFileDirPath := filepath.Join(cacheDir, logFileDir)
	err = os.MkdirAll(logFileDirPath, 0o755)
	if err != nil {
		return fmt.Errorf("make \"%s\" dir: %w", logFileDirPath, err)
	}

	logFilePath := filepath.Join(cacheDir, l.logFilePath)

	l.logFile, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o755)
	if err != nil {
		return fmt.Errorf("failed to open the file: %w", err)
	}

	return nil
}

func (l *FtLogger) SetDebugMode(value string) {
	if value == "true" {
		l.debugger = true
	}
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

func (l *FtLogger) SetLevel(level string) error {
	if level == "" {
		level = "debug"
	}
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("failed to parse level: %w", err)
	}

	l.logger.SetLevel(lvl)
	return nil
}

func (l *FtLogger) SetLoggerOutput(path string) error {
	if len(path) == 0 {
		return nil
	}
	l.logFilePath = path
	if l.isWriteToFile() {
		err := l.prepLogFile()
		if err != nil {
			return err
		}
	}
	logWriters := []io.Writer{l.logFile}
	if l.debugger {
		logWriters = append(logWriters, os.Stdout)
	}
	l.logger.SetOutput(io.MultiWriter(logWriters...))

	return nil
}

func (l *FtLogger) Close() error {
	err := l.SendLogsToController()
	if err != nil {
		logrus.Errorln(err)
	}
	if l.isWriteToFile() {
		err := l.logFile.Close()
		if err != nil {
			msg := fmt.Sprintf("close log file: %s", err)
			logrus.Error(msg)

			return errors.New(msg)
		}
	}

	return nil
}

func (l *FtLogger) SendLogsToController() error {
	if l.isWriteToFile() && l.sdr != nil {
		err := l.sendLogs()
		if err != nil {
			msg := fmt.Sprintf("send logs to controller: %s", err)
			logrus.Error(msg)

			return errors.New(msg)
		}
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return fmt.Errorf("get user cache dir: %w", err)
		}
		path := filepath.Join(cacheDir, l.logFilePath)
		err = os.Truncate(path, 0)
		if err != nil {
			logrus.Error("truncating log file: ", err)
		}
	}
	return nil
}

func (l *FtLogger) sendLogs() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get user cache dir: %w", err)
	}
	logFilePath := filepath.Join(cacheDir, l.logFilePath)
	content, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		return fmt.Errorf("error in reading log file: %w", err)
	}
	reqData := request{
		MachineID:  l.machineID(),
		LogContent: string(content),
		Instance:   l.instance,
	}

	requestBody, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("request body: %w", err)
	}

	return l.sdr.send(requestBody)
}

func (l *FtLogger) SetSender(accessToken, url, machinePairID string, machineID string) {
	l.sdr = newSender(accessToken, url, machinePairID, machineID)
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
		machineID = l.sdr.machineID
	}

	return machineID
}
