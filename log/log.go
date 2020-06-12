package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var currentTime = time.Now().Format("2006-01-02 15:04:05.000")

var (
	Trace   = &Logger{Date: currentTime, Level: "-1", Local: os.Args[0], Name: "TRACE", Path: getFileName()}
	Debug   = &Logger{Date: currentTime, Level: "0", Local: os.Args[0], Name: "DEBUG", Path: getFileName()}
	Info    = &Logger{Date: currentTime, Level: "1", Local: os.Args[0], Name: "INFO", Path: getFileName()}
	Warning = &Logger{Date: currentTime, Level: "2", Local: os.Args[0], Name: "WARN", Path: getFileName()}
	Error   = &Logger{Date: currentTime, Level: "3", Local: os.Args[0], Name: "ERROR", Path: getFileName()}
	Fatal   = &Logger{Date: currentTime, Level: "4", Local: os.Args[0], Name: "FATAL", Path: getFileName()}
	Panic   = &Logger{Date: currentTime, Level: "5", Local: os.Args[0], Name: "PANIC", Path: getFileName()}
)

// Logger estructura a Loggear
type Logger struct {
	// Date es la fecha del Log
	Date string

	// Level cuanto más alta más relevante el Log
	Level string

	// Local nombre del ejecutable
	Local string

	// Nombre del Log
	Name string

	// Path de la ubicación del Log
	Path string

	// Message principal del Log
	Message string
}

// Benchmark imprime el tiempo que transcurrio en el logger Trace.
// Ejemplo:
//  defer Benchmarkf("paso el %s","tiempo")
//  imprime: 2019/06/11 17:38:21 log.go:22: paso el tiempo: 1.2121ms
func Benchmark(fmtt string, args ...string) func() {
	started := time.Now()
	return func() {
		Trace.Pipeline(fmt.Sprintf(fmtt, args) + ": " + string(time.Since(started)))
	}
}

// Pipeline ..
func (l *Logger) Pipeline(m string) {
	var s []string
	l.Message = m
	e := reflect.ValueOf(l).Elem()

	for i := 0; i < e.NumField(); i++ {
		s = append(s, e.Field(i).String())
	}

	fmt.Println(strings.Join(s, " | "))
}

// JSON ..
func (l *Logger) JSON(m string) {
	l.Message = m
	e, _ := json.Marshal(l)
	fmt.Println(string(e))
}

func getFileName() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}

	return file + ":" + strconv.Itoa(line)
}
