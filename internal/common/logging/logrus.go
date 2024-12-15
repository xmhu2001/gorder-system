package logging

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"strings"
)

// 配置logrus格式
func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&ColorLogFormatter{})
}

type ColorLogFormatter struct{}

// Format 实现logrus.Formatter接口
func (f *ColorLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor *color.Color

	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = color.New(color.FgHiWhite)
	case logrus.InfoLevel:
		levelColor = color.New(color.FgWhite)
	case logrus.WarnLevel:
		levelColor = color.New(color.FgYellow)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = color.New(color.FgRed)
	default:
		levelColor = color.New(color.FgWhite)
	}

	levelText := levelColor.Sprintf("%s", entry.Level.String())

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	message := entry.Message

	// 处理并格式化字段
	var fields []string
	for key, value := range entry.Data {
		fields = append(fields, fmt.Sprintf("%s=%v", key, value))
	}
	fieldsStr := ""
	if len(fields) > 0 {
		fieldsStr = " " + strings.Join(fields, " ")
	}

	// 格式化日志输出
	formatted := color.New(color.FgWhite).Sprintf("%s [%s] %s%s\n", timestamp, levelText, message, fieldsStr)

	return []byte(formatted), nil
}
