package errores

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrores(t *testing.T) {
	errores := PedidoIncorrecto("detalle", nil)
	if errores.Detail != "detalle" || errores.ListaErrores != nil {
		t.Errorf("Fallo la funcion PedidoIncorrecto")
	}
}

func TestEnlistarErrores(t *testing.T) {
	data := []byte(`{"numero":true}`)
	var errJson errorJSon
	err := json.Unmarshal(data, &errJson)
	errVal := EnlistarErrores(err)
	fmt.Println(err)
	assert.Equal(t, "expected:int actual:bool", errVal[0].Message)
}

type errorJSon struct {
	Numero int `json:"numero"`
}
