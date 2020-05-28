package spew

import (
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

func Debugf(format string, v ...interface{}) {
	color.Green(format, v...)
}

func Dump(v ...interface{}) {
	s := spew.Sdump(v...)
	file, no, funcName := getCaller(2)
	color.Magenta("file: %s line: %d, funcname: %s message: %s", file, no, funcName, s)
}

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
