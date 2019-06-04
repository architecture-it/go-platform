package mq

import (
	"context"
	"gopkg.in/resty.v1"
	"fmt"
	"errors"
	"net/http"
	"time"
)

//Queue representa el objeto Queue de MQ
type Queue struct {
	name string
	api string
}

//GetQueue obtiene una cola determinada por la config especificada
func GetQueue(config Config) Queue {
	return Queue{config.QueueName,config.HTTPMQAPIUrl}
}

//Put pone un mensaje en la cola
func (q Queue) Put(data string) error {

	_,err := resty.R().
		SetBody(data).
		Put(q.api)
	return err

}

// Listen registra un callback que se invoca cada vez que llega un mensaje a la cola.
func (q Queue) Listen(ctx context.Context, f func (data string)) {
	
	go func() {
		url := fmt.Sprintf("%s/queues/%s",q.api,q.name)
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
					return
			case <-tick.C:
					if resp,err := resty.R().Get(url);err == nil && resp.StatusCode() == http.StatusOK {
						f(resp.String())
					}
			}
		}
	}()
		
}

// Publish publica en el topic 'topic' el mensaje 'data'
func Publish(topic string, data string) error {
	return errors.New("not implemented")
}