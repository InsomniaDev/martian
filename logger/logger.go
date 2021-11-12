package logger

import (
	"go/build"
	"os"

	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Logger {
	// create the logger
	logger := logrus.New()
	// with Json Formatter
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(os.Stdout)

	LOG_FILE := build.Default.GOPATH + "/config/application.log"

	file, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	logger.SetOutput(file)
	return logger
}
