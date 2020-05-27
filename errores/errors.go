package errores

import (
	"encoding/json"
	"strings"

	"github.com/go-playground/validator/v10"
)

//ErrorMessage estructura generica que maneja la compañia para los errores
type ErrorMessage struct {
	Type         string                `json:"type"`
	Title        string                `json:"title"`
	Detail       string                `json:"detail"`
	Status       int                   `json:"status"`
	ListaErrores []CampoYMensajeDeError `json:"errors"`
}

//CampoYMensajeDeError exportable para desarrollos a medida 
type CampoYMensajeDeError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

//ErrorResponse para errores no contemplados, tratar de usar los siguientes:
//PedidoIncorrecto(),RecursoNoEncontrado(),ServicioNoDisponible() para manejar 
//un estandar de errores
func ErrorResponse(titulo string, detalle string, estado int, listaDeErrores []CampoYMensajeDeError) ErrorMessage {
	return ErrorMessage{Type: "about:blank", Title: titulo, Detail: detalle, Status: estado, ListaErrores: listaDeErrores}
}

func PedidoIncorrecto(detalle string, listaDeErrores []CampoYMensajeDeError) ErrorMessage {
	return ErrorMessage{Type: "about:blank", Title: "Error en la validacion de su pedido", Detail: detalle, Status: 400, ListaErrores: listaDeErrores}
}

func RecursoNoEncontrado(detalle string, listaDeErrores []CampoYMensajeDeError) ErrorMessage {
	return ErrorMessage{Type: "about:blank", Title: "Recurso no encontrado", Detail: detalle, Status: 404, ListaErrores: listaDeErrores}
}

func ServicioNoDisponible(detalle string, listaDeErrores []CampoYMensajeDeError) ErrorMessage {
	return ErrorMessage{Type: "about:blank", Title: "Servicio no disponible momentaneamete, intente nuevamente", Detail: detalle, Status: 503, ListaErrores: listaDeErrores}
}

//EnlistarErrores decodifica el error a nuestra estructura de campo-mensaje
func EnlistarErrores(e error) []CampoYMensajeDeError {
	var unError CampoYMensajeDeError
	var ListaDeErrores []CampoYMensajeDeError
	if ute, ok := e.(*json.UnmarshalTypeError); ok {
		unError.Field = strings.ToLower(ute.Field)
		unError.Message = ute.Value
		ListaDeErrores = append(ListaDeErrores, unError)
	}

	if erroresDelValidador, ok := e.(validator.ValidationErrors); ok {
		for _, err := range erroresDelValidador {
			unError.Field = strings.ToLower(err.Field())
			if err.Param() != "" {
				unError.Message = err.Tag() + ": " + err.Param()
			} else {
				unError.Message = err.Tag()
			}
			ListaDeErrores = append(ListaDeErrores, unError)
		}
	}

	return ListaDeErrores
}
