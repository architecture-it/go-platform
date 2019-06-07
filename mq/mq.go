package mq

import (
	"context"
	"gopkg.in/resty.v1"
	"fmt"
	"strings"
	"net/http"
	"time"
	"gopkg.in/eapache/go-resiliency.v1/breaker"
)

//Queue representa el objeto Queue de MQ
type Queue struct {
	name string
	api string
}

type Topic struct {
	name string
	api string
}

var circuitBreaker *breaker.Breaker

func init() {
	circuitBreaker = breaker.New(3,1,5*time.Second)
}

func GetTopic(topic string, config Config) Topic {
	return Topic{strings.Replace(topic,"/",".",-1),config.HTTPMQAPIUrl}
}

//GetQueue obtiene una cola determinada por la config especificada
func GetQueue(config Config) Queue {
	return Queue{config.QueueName,config.HTTPMQAPIUrl}
}

//Put pone un mensaje en la cola
func (q Queue) Put(data string) error {

	url := fmt.Sprintf("%s/queues/%s",q.api,q.name)
	return circuitBreaker.Run(func() error {
		_,err := resty.R().
			SetBody(data).
			Put(url)
		return err
	})

}

// Listen registra un callback que se invoca cada vez que llega un mensaje a la cola.
func (q Queue) Listen(ctx context.Context, f func (data string)) {
	
	go func() {
		url := fmt.Sprintf("%s/queues/%s",q.api,q.name)
		tick := time.NewTicker(50*time.Millisecond)
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
func (t Topic) Publish(data string) error {

	url := fmt.Sprintf("%s/topics/%s",t.api,t.name)
	return circuitBreaker.Run(func() error {
		_,err := resty.R().
			SetBody(data).
			Post(url)
		return err
	})
}