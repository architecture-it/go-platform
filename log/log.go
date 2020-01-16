package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Fatal   *log.Logger
)

//Benchmarkf imprime el tiempo que transcurrio en el logger Trace.
//Ejemplo:
// defer Benchmarkf("paso el %s","tiempo")
// imprime: 2019/06/11 17:38:21 log.go:22: paso el tiempo: 1.2121ms
func Benchmarkf(fmtt string, args ...string) func() {
	started := time.Now()
	return func() {
		Trace.Printf("%s: %s", fmt.Sprintf(fmtt, args), time.Since(started))
	}
}
func init() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	currentTime := time.Now()

	Trace = log.New(os.Stdout,
		currentTime.Format("2006-01-02 15:04:05.000")+" | 0 | TRACE | "+filepath.Dir(d)+" | "+string(log.Llongfile)+" | ", 0)

	Info = log.New(os.Stdout,
		currentTime.Format("2006-01-02 15:04:05.000")+" | 0 | INFO | "+filepath.Dir(d)+" | "+string(log.Llongfile)+" | ", 0)

	Warning = log.New(os.Stdout,
		currentTime.Format("2006-01-02 15:04:05.000")+" | 0 | WARNING | "+filepath.Dir(d)+" | "+string(log.Llongfile)+" | ", 0)


	Error = log.New(os.Stdout,
		currentTime.Format("2006-01-02 15:04:05.000")+" | 0 | ERROR | "+filepath.Dir(d)+" | "+string(log.Llongfile)+" | ", 0)


	Fatal = log.New(os.Stdout,
		currentTime.Format("2006-01-02 15:04:05.000")+" | 0 | FATAL | "+filepath.Dir(d)+" | "+string(log.Llongfile)+" | ", 0)


}
