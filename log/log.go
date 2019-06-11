package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	Trace = log.New(ioutil.Discard,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Fatal = log.New(os.Stdout,
		"FATAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)

}
