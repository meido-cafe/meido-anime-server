package logger

import (
	"log"
	"os"
)

const (
	error_flag = log.Ldate | log.Ltime | log.Lshortfile
	info_flag  = log.Ldate | log.Ltime
	pre_info   = "[INFO] "
	pre_error  = "[ERROR] "
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, pre_info, info_flag)
	errorLogger = log.New(os.Stdout, pre_error, error_flag)
}

func Infof(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

func Info(v ...any) {
	infoLogger.Println(v)
}

func Errorf(format string, v ...any) {
	errorLogger.Printf(format, v...)
}

func Error(v ...any) {
	errorLogger.Println(v...)
}
