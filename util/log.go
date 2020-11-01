package util

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"path"
	"time"
)

/**
 * @Description: logrus日志库,设置log等级
 * @param level
 */
func setLogLevel(level string) {
	switch level {
	case "InfoLevel":
		logrus.SetLevel(logrus.InfoLevel)
	case "DebugLevel":
		logrus.SetLevel(logrus.DebugLevel)
	case "ErrorLevel":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FatalLevel":
		logrus.SetLevel(logrus.FatalLevel)
	case "PanicLevel":
		logrus.SetLevel(logrus.PanicLevel)
	case "TraceLevel":
		logrus.SetLevel(logrus.TraceLevel)
	case "WarnLevel":
		logrus.SetLevel(logrus.WarnLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

/**
 * @Description: 设置log格式
 * @param formatter
 */
func setLogFormatter(formatter string) {
	if formatter == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
}

/**
 * @Description: 设置log输出文件
 */
func setLogOutput() {
	now := time.Now()
	//获取日志文件路径
	logFilePath := viper.GetString("log.filePath")
	if logFilePath == "" {
		if dir, err := os.Getwd(); err == nil {
			logFilePath = dir + "/logs/"
		}
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			log.Printf(err.Error())
		}
	}

	//设置日志文件名字
	logFileName := viper.GetString("log.fileName")
	//当前日期时间命名
	if logFileName == "" {
		logFileName = now.Format("2006-01-02") + ".log"
	}
	//日志文件,拼接两部分
	fileName := path.Join(logFilePath, logFileName)

	//写入文件
	var writer io.Writer
	//获取配置文件中的日志文件类型
	logFileType := viper.GetString("log.fileType")
	//未设置和Stdout设置为标准输出
	if logFileType == "File" {
		var err error
		writer, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Panic("create file log.txt failed: ", err)
		}
	} else {
		writer = os.Stdout
	}

	logrus.SetOutput(io.MultiWriter(writer))

}

func InitLog() {
	//设置输出
	setLogOutput()
	//设置日志级别
	setLogLevel(viper.GetString("log.level"))
	//设置日志格式
	setLogFormatter(viper.GetString("log.formatter"))

}
