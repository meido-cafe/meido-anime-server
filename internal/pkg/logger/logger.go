package logger

import (
	"log"
	"os"
)

const (
	flag      = log.Ldate | log.Ltime | log.Lshortfile
	pre_info  = "[INFO] "
	pre_error = "[ERROR] "
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, pre_info, flag)
	errorLogger = log.New(os.Stdout, pre_error, flag)
}

func Infof(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

func Info(v ...interface{}) {
	infoLogger.Println(v)
}

func Errorf(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func Error(v ...interface{}) {
	errorLogger.Println(v...)
}
