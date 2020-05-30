package log

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	hostName, _ = os.Hostname()

	debugWriter, normalWriter io.Writer

	glog = logrus.New()

	// for purger
	Log = glog.WithFields(logrus.Fields{
		"role": "hke-control",
	})
)

const NullLogPrefix = "$$$$$$$$$$$$$$$$$$$$$$$$$$$$"

type Config struct {
	Console  bool   `toml:"console"`
	Type     string `toml:"type"`
	Level    string `toml:"level"`
	FileName string `toml:"filename"`
	Path     string `toml:"path"`
	Buffer   int    `toml:"buffer"`
	MaxAge   int    `toml:"maxage"`
	Rotation int    `toml:"rotation"`
}

func InitLogger(c *Config) func() {
	formatter := &DefaultFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HostName:        hostName,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", findCaller()
		},
		// CallerPrettyfier: func(f *runtime.Frame) (string, string) {
		// 	fileName := f.File[strings.Index(f.File, "pusher")+7:]
		// 	return "", fmt.Sprintf("%s:%d", fileName, f.Line)
		// },
	}

	glog.SetOutput(ioutil.Discard)
	glog.SetFormatter(
		// default std
		&logrus.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05.000",
		},
		// formatter,
	)

	// open console stdout
	if c.Console {
		glog.SetOutput(os.Stdout)
	}

	var err error
	if debugWriter, err = rotatelogs.New(
		filepath.Join(c.Path, "debug.log.%Y%m%d"),
		rotatelogs.WithLinkName(filepath.Join(c.Path, "debug.log")),
		rotatelogs.WithMaxAge(time.Duration(c.MaxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(c.Rotation)*time.Hour),
	); err != nil {
		panic(err)
	}

	if normalWriter, err = rotatelogs.New(
		filepath.Join(c.Path, c.FileName+".%Y%m%d"),
		rotatelogs.WithLinkName(filepath.Join(c.Path, c.FileName)),
		rotatelogs.WithMaxAge(time.Duration(c.MaxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(c.Rotation)*time.Hour),
	); err != nil {
		panic(err)
	}

	if c.Buffer > 0 {
		debugWriter = bufio.NewWriterSize(debugWriter, c.Buffer)
		normalWriter = bufio.NewWriterSize(normalWriter, c.Buffer)
	}

	lfHook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: normalWriter,
			logrus.InfoLevel:  normalWriter,
			logrus.WarnLevel:  normalWriter,
			logrus.ErrorLevel: normalWriter,
			logrus.FatalLevel: normalWriter,
			logrus.PanicLevel: normalWriter,
		},
		&logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05"},
	)

	switch c.Type {
	case "json":
		lfHook.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		lfHook.SetFormatter(&logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		lfHook.SetFormatter(formatter)
	}

	switch c.Level {
	case "debug":
		glog.SetLevel(logrus.DebugLevel)
	case "info":
		glog.SetLevel(logrus.InfoLevel)
	case "warn":
		glog.SetLevel(logrus.WarnLevel)
	case "error":
		glog.SetLevel(logrus.ErrorLevel)
	case "fatal":
		glog.SetLevel(logrus.FatalLevel)
	case "panic":
		glog.SetLevel(logrus.PanicLevel)
	default:
		glog.SetLevel(logrus.InfoLevel)
	}

	glog.AddHook(lfHook)

	return Flush
}

// only used for buffer log
func Flush() {
	if buffWriter, ok := debugWriter.(*bufio.Writer); ok {
		buffWriter.Flush()
	}

	if buffWriter, ok := normalWriter.(*bufio.Writer); ok {
		buffWriter.Flush()
	}
}

// Formatter implements logrus.Formatter interface
type DefaultFormatter struct {
	TimestampFormat  string
	HostName         string
	CallerPrettyfier func(f *runtime.Frame) (string, string)
}

// Format building log message
func (f *DefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// time field
	b.WriteString(entry.Time.Format(f.TimestampFormat))

	// hostname field
	b.WriteString("$$" + f.HostName)

	// level field
	b.WriteString("$$" + strings.ToUpper(entry.Level.String()))

	// component field
	b.WriteString("$$" + entry.Data["role"].(string))

	_, fileVal := f.CallerPrettyfier(entry.Caller)
	// file no
	b.WriteString("$$" + fileVal)

	// msg field
	b.WriteString("$$" + entry.Message)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Printf(format string, args ...interface{}) {
	Log.Printf(format, args...)
}

func Print(args ...interface{}) {
	Log.Print(args...)
}

func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func findCaller() string {
	var (
		funcName = ""
		file     = ""
		line     = 0
		pc       uintptr
	)

	// logrus + lfshook + log.go = 12
	for i := 12; i < 15; i++ {
		file, line, pc = getCaller(i)
		// fileter logrus + lfshook + log.go
		if strings.HasPrefix(file, "log/log.go") {
			continue
		}
		if strings.HasPrefix(file, "logrus") {
			continue
		}
		if strings.Contains(file, "lfshook.go") {
			continue
		}
		break
	}

	fullFnName := runtime.FuncForPC(pc)

	if fullFnName != nil {
		fnNameStr := fullFnName.Name()
		parts := strings.Split(fnNameStr, ".")
		funcName = parts[len(parts)-1]
	}

	return fmt.Sprintf("%s:%d:%s()", file, line, funcName)
}

func getCaller(skip int) (string, int, uintptr) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, pc
	}

	n := 0

	// get package name
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line, pc
}

func hanlePanicf(format string, args ...interface{}) {
}
