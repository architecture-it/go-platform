package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// log.Info.Println("Info")
//

// Benchmark - Imprime el tiempo que transcurrio en el logger Trace.
// Ejemplo:
//  defer Benchmarkf("paso el %s","tiempo")
//  imprime: 2019/06/11 17:38:21 log.go:22: paso el tiempo: 1.2121ms
func Benchmark(fmtt string, args ...string) func() {
	//started := time.Now()
	return func() {
		//Trace.Printf("%s: %s", fmt.Sprintf(fmtt, args), time.Since(started))
	}
}

func main() {
	// currentTime := time.Now()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log := zerolog.New(output).With().Timestamp().Logger()

	log.Info().Msg("Hello World")
}
