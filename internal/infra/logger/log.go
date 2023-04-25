package logger

import (
	"github.com/sirupsen/logrus"
)

func Init(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{
		logrus.FieldKeyMsg:  "message",
		logrus.FieldKeyTime: "timestamp",
	}})
}
