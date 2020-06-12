package errores

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrores(t *testing.T) {
	errores := PedidoIncorrecto.Default("detalle", nil)
	if errores.Detail != "detalle" || errores.List != nil {
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

	assert.Equal(t, string(out), `[{"name":"numero","description":"bool"},{"name":"numero","description":"string"}]`)
}

type errJSON struct {
	Numero int `json:"numero"`
}

func TestDefault(t *testing.T) {
	var errJSON errJSON

	err1 := json.Unmarshal([]byte(`{"numero":true}`), &errJSON)
	err2 := json.Unmarshal([]byte(`{"numero":1}`), &errJSON)

	errVal := ErrorResponse.Default("Default", err1, err2)

	assert.Equal(t, errVal.List[0].Description, "bool")
}
