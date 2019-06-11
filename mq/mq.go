package mq

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/eapache/go-resiliency.v1/breaker"
	"gopkg.in/resty.v1"
)

//Queue representa el objeto Queue de MQ
type Queue struct {
	name string
	api  string
}

//Topic representa al objeto Topic de MQ
type Topic struct {
	name string
	api  string
}

var circuitBreaker *breaker.Breaker

func init() {
	circuitBreaker = breaker.New(3, 1, 5*time.Second)
}

//GetTopic obtiene un topic para publicar.
func GetTopic(topic string, config Config) Topic {
	return Topic{strings.Replace(topic, "/", ".", -1), config.HTTPMQAPIUrl}
}

//GetQueue obtiene una cola determinada por la config especificada
func GetQueue(config Config) Queue {
	return Queue{config.QueueName, config.HTTPMQAPIUrl}
}

//Put pone un mensaje en la cola
func (q Queue) Put(data string) error {

	url := fmt.Sprintf("%s/queues/%s", q.api, q.name)
	return circuitBreaker.Run(func() error {
		_, err := resty.R().
			SetBody(data).
			Put(url)
		return err
	})

}

// Listen registra un callback que se invoca cada vez que llega un mensaje a la cola.
func (q Queue) Listen(ctx context.Context, f func(data string)) {

	go func() {
		url := fmt.Sprintf("%s/queues/%s", q.api, q.name)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if resp, err := callMqBridge(url); err == nil {
					f(resp)
				}
			}
		}
	}()

}

func callMqBridge(url string) (string, error) {
	//defer log.Benchmarkf("call resty: ", url)
	resp, err := resty.R().Get(url)
	if err == nil && resp.StatusCode() == http.StatusOK {
		return resp.String(), nil
	}
	return "", err
}

// Publish publica en el topic 'topic' el mensaje 'data'
func (t Topic) Publish(data string) error {

	url := fmt.Sprintf("%s/topics/%s", t.api, t.name)
	return circuitBreaker.Run(func() error {
		_, err := resty.R().
			SetBody(data).
			Post(url)
		return err
	})
}
