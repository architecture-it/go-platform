package health

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/architecture-it/go-platform/log"
)

type IbmMQStatus struct {
	Status string `json:"status"`
}

func IbmMQHealthChecker() Checker {
	checker := Checker{Health: Health{Status: Status{Code: UP, Description: ""}, Details: ""}, Name: "ibmQueueHealthIndicator"}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(os.Getenv("HTTP_MQ_API_URL") + "/health")
	if err != nil {
		checker.Health.Status.Code = DOWN
		log.SugarLogger.Error("Error al obtener el health del MQ. ", err)
		return checker
	}

	if resp.StatusCode == http.StatusOK {
		jsonMqHealth, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.SugarLogger.Error("Error al leer la respuesta del health de la API de MQ.", err)
			checker.Health.Status.Code = DOWN
			return checker
		}
		var result IbmMQStatus
		json.Unmarshal(jsonMqHealth, &result)
		checker.Health.Status.Code = result.Status
	} else {
		checker.Health.Status.Code = DOWN
	}
	return checker
}
