package errores

import (
	"encoding/json"
	"strings"

	"github.com/go-playground/validator"
)

var (
	PedidoIncorrecto     ErrorResponse
	RecursoNoEncontrado  ErrorResponse
	ErrorServidorInterno ErrorResponse
	ServicioNoDisponible ErrorResponse
)

// ErrorResponse Estructura genérica para errores de Andreani S.A.
type ErrorResponse struct {
	Type   string  `json:"type" default1:"about:blank"`
	Title  string  `json:"title"`
	Detail string  `json:"detail"`
	Status int     `json:"status"`
	Errors []Error `json:"errors"`
}

// Error Exportable para desarrollos a medida
type Error struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func init() {
	PedidoIncorrecto = ErrorResponse{Type: "about:blank", Title: "Error en la validacion de su pedido", Status: 400}
	RecursoNoEncontrado = ErrorResponse{Type: "about:blank", Title: "Recurso no encontrado", Status: 404}
	ErrorServidorInterno = ErrorResponse{Type: "about:blank", Title: "Error en la Respuesta", Status: 500}
	ServicioNoDisponible = ErrorResponse{Type: "about:blank", Title: "Servicio no disponible momentaneamete, intente nuevamente", Status: 503}
}

// Default - Permite settear solo el detalle y los errores del campo List
func (er ErrorResponse) Default(detalle string, errores ...error) (int, *ErrorResponse) {
	per := &er
	per.Detail = detalle
	per.Errors = *errores2List(errores)
	return er.Status, per
}

// All - Permite settear todos los valores del error
func (er ErrorResponse) All(tipo string, titulo string, detalle string, status int, errores ...error) (int, *ErrorResponse) {
	per := &er
	per.Type = tipo
	per.Title = titulo
	per.Detail = detalle
	per.Status = status
	per.Errors = *errores2List(errores)
	return er.Status, per
}

func errores2List(errs []error) *[]Error {
	var fieldList []Error
	for _, err := range errs {
		var field Error
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			field.Name = strings.ToLower(ute.Field)
			field.Message = ute.Value
			fieldList = append(fieldList, field)
		}

		if validatorErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validatorErrors {
				field.Name = strings.ToLower(e.Field())
				if e.Param() != "" {
					field.Message = e.Tag() + ": " + e.Param()
				} else {
					field.Message = e.Tag()
				}
				fieldList = append(fieldList, field)
			}
		}
	}

	return &fieldList

}
