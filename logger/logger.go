package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	Logger.controllerURL, ok = os.LookupEnv("CONTROLLER_URL")
	if !ok {
		Logger.controllerURL = "https://dev-controller.gather.network"
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
	Logger        *logrus.Entry
	logFile       *os.File
	logFilePath   string
	controllerURL string
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

type SendLogsRequest struct {
	Instance   string `json:"instance"`
	LogContent string `json:"log_content"`
	MachineID  string `json:"machine_id"`
}

func (l *FtLogger) SendLogs(auth string, machinePairID string) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		Logger.Fatalf("get user cache dir: %s", err)
		return err
	}
	logFilePath := filepath.Join(cacheDir, l.logFilePath)
	content, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		Logger.Fatalf("error in reading log file: ", err)
	}
	machineID := ""
	if instance == "sender" || instance == "receiver" {
		machineID = "2"
	}

	requestBody, err := json.Marshal(SendLogsRequest{instance, string(content), machineID})
	if err != nil {
		Logger.Errorf("request body: ", err)
		return err
	}
	req, err := http.NewRequest("POST", Logger.controllerURL+"/api/v1/file_transfer/machine_pair/"+machinePairID+"/logs/", bytes.NewBuffer(requestBody))
	if err != nil {
		Logger.Errorf("creating request: ", err)
		return err
	}

	req.Header.Add("Authorization-Token", auth)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Logger.Errorf("sending request:", err)
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Errorf("unexpected error reading response:", err)
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(bodyBytes), &data)
	if err != nil {
		Logger.Errorf("unexpected error in unmarshal:", err)
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		Logger.Infof("Succes response: ", data["details"])
		return nil

	case http.StatusBadRequest:
		if msg, ok := data["instance"]; ok {
			Logger.Errorf("Failed: ", msg)
			return fmt.Errorf("%v", msg)
		}
		if msg, ok := data["machine_id"]; ok {
			Logger.Errorf("Failed: ", msg)
			return fmt.Errorf("%v", msg)
		}
		return fmt.Errorf("Error: 400")

	case http.StatusNotFound:
		if msg, ok := data["detail"]; ok {
			Logger.Errorf("Failed: ", msg)
			return fmt.Errorf("%v", msg)
		}
		return fmt.Errorf("Error: 404")

	case http.StatusUnauthorized:
		if msg, ok := data["error"]; ok {
			Logger.Errorf("Failed: ", msg)
			return fmt.Errorf("%v", msg)
		}
		return fmt.Errorf("Error: 401")
	}
	return fmt.Errorf("Error Response: %v", resp.StatusCode)
}
