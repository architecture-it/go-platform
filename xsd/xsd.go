package xsd

import (
	"encoding/xml"
	"errors"
	"reflect"
	"strings"
	"time"
)

var namespaces map[string]string

func init() {
	namespaces = map[string]string{
		"ei": "http://integracion.andreani.com/eventosDeIntegracion",
		"dr": "http://integracion.andreani.com/datosDeReferencia/",
		"ea": "http://integracion.andreani.com/eventosDeAlmacenes/",
		"di": "http://integracion.andreani.com/discovery",
		"in": "http://integracion.andreani.com/incidencias",
		"pr": "http://integracion.andreani.com/preguntas",
		"re": "http://integracion.andreani.com/respuestas",
	}

}

//GetEvento devuelve el xml del evento a partir de la estructura. Permite la opcion de devolver el evento indentado
func GetEvento(d interface{}, indent ...bool) (string, error) {

	field, valid := reflect.TypeOf(d).Elem().FieldByName("XMLNs")
	if valid != true {
		return "", errors.New("No se pudo acceder al campo XMLNs de la estructura")
	}

	for alias, ns := range namespaces {
		match := strings.Contains(string(field.Tag), alias)
		if match {
			b, err := getXML(d, ns, indent[0])
			return string(b), err
		}
	}
	return "", errors.New("No se pudo encontrar un namespace para la estructura")

}

//GetDate debe ser usada siempre que se quiera setear un Date en la estructura
func GetDate(b time.Time) *xsdDate {
	return (*xsdDate)(&b)
}

//GetTime debe ser usada siempre que se quiera setear un Time en la estructura
func GetTime(b time.Time) *xsdTime {
	return (*xsdTime)(&b)
}

func getXML(d interface{}, namespace string, indent bool) ([]byte, error) {
	v := reflect.ValueOf(d).Elem()

	if f := v.FieldByName("XMLNs"); f.IsValid() {
		f.SetString(f.Interface().(string) + namespace)
	} else {
		return nil, errors.New("Error al acceder al campo XMLNs de la estructura")
	}
	if f := v.FieldByName("Timestamp"); f.IsValid() {
		f.Set(reflect.ValueOf(GetDate(time.Now())))
	}

	if indent {
		return xml.MarshalIndent(d, "", " ")
	}
	return xml.Marshal(d)

}
