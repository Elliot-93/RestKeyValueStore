package logger

import (
	"fmt"
	"log"
	"os"
)

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

func Info(msg interface{}) {
	logToStdOutAndFile("INFO ", msg, StandardLogFile)
}

func Warning(msg interface{}) {
	logToStdOutAndFile("WARNING ", msg, StandardLogFile)
}

func Error(msg interface{}) {
	logToStdOutAndFile("ERROR ", msg, StandardLogFile)
}

func Fatal(msg interface{}) {
	logToStdOutAndFile("FATAL ", msg, StandardLogFile)
}

func LogRequest(requestDetails interface{}) {
	logToStdOutAndFile("INFO ", requestDetails, RequestsLogFile)
}

func formatWithLogLevel(prefix, msg interface{}) string {
	return fmt.Sprintf("%s: %v", prefix, msg)
}

func logToStdOutAndFile(prefix string, msg interface{}, logFile string) {
	logToStdOut(formatWithLogLevel(prefix, msg))
	logToFile(prefix, msg, logFile)
}

func logToStdOut(msg interface{}) {
	fmt.Println(msg)
}

func logToFile(prefix string, message interface{}, logfile string) {
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger := log.New(file, prefix, log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println(message)
}
