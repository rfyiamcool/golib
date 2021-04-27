package spews

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
)

// usage set fg color and bg color
// red := color.New(color.FgRed).Add(color.BgGreen)
// red.Println("xiaorui.cc")

func Open() {
	color.Set(color.FgMagenta, color.Bold)
}

func Close() {
	defer color.Unset()
}

// usage set fg color and bg color
// red := color.New(color.FgRed).Add(color.BgGreen)
// red.Println("xiaorui.cc")

func Panic(v ...interface{}) {
	red := color.New(color.FgBlack).Add(color.BgRed)
	red.Println(v...)
}

func Error(v ...interface{}) {
	s := fmt.Sprintln(v...)
	color.Red(s)
}

func Errorf(format string, v ...interface{}) {
	color.Red(format, v...)
}

func Warn(v ...interface{}) {
	s := fmt.Sprintln(v...)
	color.Yellow(s)
}

func Warnf(format string, v ...interface{}) {
	color.Yellow(format, v...)
}

func Info(v ...interface{}) {
	s := fmt.Sprintln(v...)
	color.Blue(s)
}

func Infof(format string, v ...interface{}) {
	color.Blue(format, v...)
}

func Debug(v ...interface{}) {
	s := fmt.Sprintln(v...)
	color.Green(s)
}

func Alert(v ...interface{}) {
	s := fmt.Sprintln(v...)
	color.HiBlue(s)
}

func Debugf(format string, v ...interface{}) {
	color.Green(format, v...)
}

func JsonDump(vlist ...interface{}) {
	file, no, funcName := getCaller(2)
	for _, v := range vlist {
		bs, _ := json.Marshal(v)
		color.Magenta("file: %s line: %d, funcname: %s, message: %s", file, no, funcName, string(bs))
	}
}

func Dump(vlist ...interface{}) {
	var values []interface{}

	for _, v := range vlist {
		switch v.(type) {
		case []byte:
			values = append(values, string(v.([]byte)))
		default:
			values = append(values, v)
		}
	}

	s := spew.Sdump(values...)
	file, no, funcName := getCaller(2)
	color.Magenta("file: %s line: %d, funcname: %s, message: %s", file, no, funcName, s)
}

func Stack(v ...interface{}) {
	stack := getStack()
	v = append(v, "stack: %s", stack)
	Debug(v...)
}

// getCaller get filename, line, fucntion name
func getCaller(skip int) (string, int, string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, ""
	}

	var (
		n        = 0
		funcName string
	)

	// get package name
	for i := len(file) - 1; i > 0; i-- {
		if file[i] != '/' {
			continue
		}
		n++
		if n >= 2 {
			file = file[i+1:]
			break
		}
	}

	fnpc := runtime.FuncForPC(pc)

	if fnpc != nil {
		fnNameStr := fnpc.Name()
		parts := strings.Split(fnNameStr, ".")
		funcName = parts[len(parts)-1]
	}

	return file, line, funcName
}

// getStack get full function stack
func getStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
