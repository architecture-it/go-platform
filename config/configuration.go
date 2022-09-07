package configuration

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func GetConfiguration(filename string) map[string]interface{} {
	var configurations map[string]interface{}
	configurations = make(map[string]interface{})

	// upload enviroment file
	if len(filename) != 0 {
		fileExtension := filepath.Ext(filename)

		if fileExtension == ".json" {
			tpl, err := template.ParseFiles(filename)
			if err != nil {
				log.Fatal(err)
			}
			var ts bytes.Buffer

			var struc interface{}

			err = tpl.Execute(&ts, struc) // Execute will fill the buffer so pass as reference
			if err != nil {
				log.Fatal(err)
			}
			json.Unmarshal(ts.Bytes(), &configurations)
		}

		if fileExtension == ".yaml" {

			bytes, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Printf("yamlFile.Get err   #%v ", err)
			}

			err = yaml.Unmarshal(bytes, &configurations)
			if err != nil {
				log.Fatalf("Unmarshal: %v", err)
			}
		}
	}

	// upload enviroment system
	for _, item := range os.Environ() {
		items := strings.Split(item, "=")
		configurations[items[0]] = items[1]
	}

	return configurations
}
