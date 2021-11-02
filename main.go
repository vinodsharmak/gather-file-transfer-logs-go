package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
	// Only log the warning severity or above.
}
func Fatal(s string) {
	logrus.WithFields(logrus.Fields{
		"instance": "sender",
	}).Fatal(s)
}
func main() {
	Fatal("hello There")
}
