package errores

import (
	"testing"
)

func TestErrores(t *testing.T) {
	errores := PedidoIncorrecto("detalle", nil)
	if errores.Detail != "detalle" || errores.ListaErrores != nil {
		t.Errorf("Fallo la funcion PedidoIncorrecto")
	}
}
