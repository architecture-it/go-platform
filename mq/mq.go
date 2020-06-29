package mq

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
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
var (
	errorQueueEmpty = errors.New("queue empty")
)

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
		res, err := resty.R().
			SetBody(data).
			Put(url)
		if res.StatusCode() != http.StatusOK {
			err = errors.New("API Bridge falla al publicar el mensaje")
		}
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
					for msg := range resp {
						f(msg)
					}
				}
			}
		}
	}()

}

func callMqBridge(url string) (<-chan string, error) {
	//defer log.Benchmarkf("call resty: ", url)
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s?batch=20", url), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)

	code := resp.StatusCode

	if err == nil && code == http.StatusOK {
		return parseMultipartResponse(resp), nil
	}
	if err == nil && code == http.StatusNoContent {
		err = errorQueueEmpty
	}
	return nil, err
}

func parseMultipartResponse(resp *http.Response) <-chan string {

	result := make(chan string, 10)

	go func() {
		_, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		mr := multipart.NewReader(resp.Body, params["boundary"])
		defer resp.Body.Close()
		defer close(result)
		for part, err := mr.NextPart(); err == nil; part, err = mr.NextPart() {
			value, _ := ioutil.ReadAll(part)
			result <- string(value)
		}
	}()
	return result
}

// Publish publica en el topic 'topic' el mensaje 'data'
func (t Topic) Publish(data string) error {

	url := fmt.Sprintf("%s/topics/%s", t.api, t.name)
	return circuitBreaker.Run(func() error {
		res, err := resty.R().
			SetBody(data).
			Post(url)
		if res.StatusCode() != http.StatusOK {
			err = errors.New("API Bridge falla al publicar el mensaje")
		}
		return err
	})
}
