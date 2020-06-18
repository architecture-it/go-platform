package errores

import (
	"encoding/json"
	"strings"

	"github.com/go-playground/validator"
)

var (
	ErrorResponse        *ErrorRequest
	PedidoIncorrecto     *ErrorRequest
	RecursoNoEncontrado  *ErrorRequest
	ServicioNoDisponible *ErrorRequest
)

// ErrorRequest Estructura genérica para errores de Andreani S.A.
type ErrorRequest struct {
	Type   string  `json:"type"`
	Title  string  `json:"title"`
	Detail string  `json:"detail"`
	Status int     `json:"status"`
	List   []Field `json:"errors"`
}

// Field Exportable para desarrollos a medida
type Field struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func init() {
	ErrorResponse = &ErrorRequest{Type: "about:blank", Title: "Error en la Respuesta", Status: 500}
	PedidoIncorrecto = &ErrorRequest{Type: "about:blank", Title: "Error en la validacion de su pedido", Status: 400}
	RecursoNoEncontrado = &ErrorRequest{Type: "about:blank", Title: "Recurso no encontrado", Status: 404}
	ServicioNoDisponible = &ErrorRequest{Type: "about:blank", Title: "Servicio no disponible momentaneamete, intente nuevamente", Status: 503}
}

// Default - Permite settear solo el detalle y los errores del campo List
func (er *ErrorRequest) Default(d string, e ...error) ErrorRequest {
	er.Detail = d
	er.List = errores2List(e)
	return *er
}

// All - Permite settear todos los valores del error
func (er *ErrorRequest) All(t string, ti string, d string, s int, e ...error) ErrorRequest {
	er.Type = t
	er.Title = ti
	er.Detail = d
	er.Status = s
	er.List = errores2List(e)
	return *er
}

func errores2List(errs []error) []Field {
	var fieldList []Field

	for _, err := range errs {
		var field Field
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			field.Name = strings.ToLower(ute.Field)
			field.Description = ute.Value
			fieldList = append(fieldList, field)
		}

		if validatorErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validatorErrors {
				field.Name = strings.ToLower(e.Field())
				if e.Param() != "" {
					field.Description = e.Tag() + ": " + e.Param()
				} else {
					field.Description = e.Tag()
				}
				fieldList = append(fieldList, field)
			}
		}
	}

	return fieldList

}
