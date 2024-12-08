package logging

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// 配置logrus格式
func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
}

func SetFormatter(logger *logrus.Logger) {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyMsg:   "message",
		},
	})
	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
		//logger.SetFormatter(&prefixed.TextFormatter{
		//	ForceFormatting: true,
		//})
	}
}
