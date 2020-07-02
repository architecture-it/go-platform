package xsd

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestOrdenDeEnvioSolicitada(t *testing.T) {
	ordenDeEnvio := OrdenDeEnvioSolicitada{Remitente: "APIs", Cuando: GetDate(time.Now())}
	str, err := GetEvento(&ordenDeEnvio, true)
	if err != nil {
		panic(err)
	}

	file := "OrdenDeEnvioSolicitada.xml"
	er := ioutil.WriteFile(file, []byte(str), 0644)
	if er != nil {
		panic(er)
	}
	err = os.Remove(file)
	if er != nil {
		panic(er)
	}
}
