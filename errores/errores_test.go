package errores

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrores(t *testing.T) {
	_, errores := PedidoIncorrecto.Default("detalle", nil)
	if errores.Detail != "detalle" || errores.Errors != nil {
		t.Errorf("Fallo la funcion PedidoIncorrecto")
	}
}

func TestErrores2List(t *testing.T) {
	var errJSON errJSON
	data1 := []byte(`{"numero":true}`)
	err1 := json.Unmarshal(data1, &errJSON)
	data2 := []byte(`{"numero":"true"}`)
	err2 := json.Unmarshal(data2, &errJSON)
	errVal := errores2List([]error{err1, err2})
	out, err := json.Marshal(errVal)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(out), `[{"name":"numero","message":"bool"},{"name":"numero","message":"string"}]`)
}

type errJSON struct {
	Numero int `json:"numero"`
}

func TestDefault(t *testing.T) {
	var errJSON errJSON

	err1 := json.Unmarshal([]byte(`{"numero":true}`), &errJSON)
	err2 := json.Unmarshal([]byte(`{"numero":1}`), &errJSON)

	_, errVal := ErrorServidorInterno.Default("Default", err1, err2)

	assert.Equal(t, errVal.Errors[0].Message, "bool")
}

func TestVariosAll(t *testing.T) {
	var errJSON errJSON

	err1 := json.Unmarshal([]byte(`{"numero":true}`), &errJSON)
	err2 := json.Unmarshal([]byte(`{"numero":1}`), &errJSON)

	ErrorServidorInterno.All("Tipo", "Titulo", "Detalle", 200, err1, err2)
	_, errVal := ErrorServidorInterno.Default("Este es un ejemplo")

	assert.Equal(t, errVal.Title, "Error en la Respuesta")

}
