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

func IbmMQHealthChecker() Status {
	health := Status{"IbmQueueHealthIndicator", StatusResult(UP), ""}

	if os.Getenv("HTTP_MQ_API_URL") == "" {
		health.Status = NOT_SET
		return health
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(os.Getenv("HTTP_MQ_API_URL") + "/health")
	if err != nil {
		health.Status = DOWN
		log.Error.Println("Error al obtener el health del MQ. ", err)
		return health
	}

	if resp.StatusCode == http.StatusOK {
		jsonMqHealth, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error.Println("Error al leer la respuesta del health de la API de MQ.", err)
			health.Status = DOWN
			return health
		}
		var result IbmMQStatus
		json.Unmarshal(jsonMqHealth, &result)
		health.Status = StatusResult(result.Status)
	} else {
		health.Status = DOWN
	}
	return health
}
