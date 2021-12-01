package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	infoPrefix    = "INFO "
	warningPrefix = "WARNING "
	errorPrefix   = "ERROR "
	fatalPrefix   = "FATAL "
)

var infoLogger *log.Logger
var warningLogger *log.Logger
var errorLogger *log.Logger
var fatalLogger *log.Logger
var requestLogger *log.Logger

type Logger interface {
	LogRequest()
	Info(msg interface{})
	Warning(msg interface{})
	Error(msg interface{})
	Fatal(msg interface{})
}

const (
	StandardLogFile = "store.log"
	RequestsLogFile = "htaccess.log"
)

func init() {
	infoLogger = buildLogger(infoPrefix, StandardLogFile)
	warningLogger = buildLogger(warningPrefix, StandardLogFile)
	errorLogger = buildLogger(errorPrefix, StandardLogFile)
	fatalLogger = buildLogger(fatalPrefix, StandardLogFile)
	requestLogger = buildLogger(infoPrefix, RequestsLogFile)
}

func Info(msg interface{}) {
	logToStdOutAndFile(infoLogger, infoPrefix, msg)
}

func Warning(msg interface{}) {
	logToStdOutAndFile(warningLogger, warningPrefix, msg)
}

func Error(msg interface{}) {
	logToStdOutAndFile(errorLogger, errorPrefix, msg)
}

func Fatal(msg interface{}) {
	logToStdOutAndFile(fatalLogger, fatalPrefix, msg)
}

func LogRequest(requestDetails interface{}) {
	logToStdOutAndFile(requestLogger, infoPrefix, requestDetails)
}

func logToStdOutAndFile(logger *log.Logger, prefix string, msg interface{}) {
	logToStdOut(formatWithLogLevel(prefix, msg))
	logger.Println(msg)
}

func formatWithLogLevel(prefix, msg interface{}) string {
	return fmt.Sprintf("%s: %v", prefix, msg)
}
func logToStdOut(msg interface{}) {
	fmt.Println(msg)
}

func buildLogger(prefix string, logfile string) *log.Logger {
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}

	return log.New(file, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}
