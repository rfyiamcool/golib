package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	ErrInvalidLevel = errors.New("log level setting error")

	defaultStdout = os.Stdout
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warnLogger    *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger
)

func init() {
	reset()
}

func reset() {
	debugLogger = log.New(defaultStdout, "[Debug] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(defaultStdout, "[Info] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger = log.New(defaultStdout, "[Wain] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(defaultStdout, "[Error] ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(defaultStdout, "[Fatal] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func SetStdoutFile(filePath string) error {
	fileFd, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defaultStdout = fileFd
	reset()
	return nil
}

func SetLevelWithDefault(lv, defaultLv string) {
	err := SetLevel(lv)
	if err != nil {
		SetLevel(defaultLv)
	}
}

func SetLevel(lv string) error {
	if lv == "" {
		return ErrInvalidLevel
	}

	var (
		l     = strings.ToUpper(lv)
		level int
	)

	switch l {
	case "DEBUG":
		level = 1
	case "INFO":
		level = 2
	case "WARN":
		level = 3
	case "ERROR":
		level = 4
	case "FATAL":
		level = 5
	default:
		level = 6
	}

	if level == 6 {
		return ErrInvalidLevel
	}

	return nil
}

func Debug(format string, v ...interface{}) {
	if 1 >= level {
		debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	if 2 >= level {
		infoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Warn(format string, v ...interface{}) {
	if 3 >= level {
		warnLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(format string, v ...interface{}) {
	if 4 >= level {
		errorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Fatal(format string, v ...interface{}) {
	if 5 >= level {
		fatalLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Traceln(v ...interface{}) {
	if 0 >= level {
		traceLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Debugln(v ...interface{}) {
	if 1 >= level {
		debugLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Infoln(v ...interface{}) {
	if 2 >= level {
		infoLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Warnln(v ...interface{}) {
	if 3 >= level {
		warnLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Errorln(v ...interface{}) {
	if 4 >= level {
		errorLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Fatalln(v ...interface{}) {
	if 5 >= level {
		fatalLogger.Output(2, fmt.Sprintln(v...))
	}
}
