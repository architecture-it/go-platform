package mq

import (
	"context"
	"gopkg.in/resty.v1"
	"fmt"
	"errors"
	"net/http"
)

type Queue struct {
	name string
	api string
}

type Listener struct {
	callback func (data string) 
	queue Queue
	context context.Context 
}

func GetQueue(config Config) Queue {
	return Queue{config.QueueName,config.HTTPMQAPIUrl}
}
func (q Queue) Put(data string) error {

	_,err := resty.R().
		SetBody(data).
		Put(q.api)
	return err

}

func (q Queue) Listen(ctx context.Context, f func (data string)) {
	
	go func() {
		exit:=false

		go func() {
			<- ctx.Done()
			exit = true
		}()

		url := fmt.Sprintf("%s/queues/%s",q.api,q.name)
		for {
			if exit == true {
				break
			}
			
			if resp,err := resty.R().Get(url);err == nil && resp.StatusCode() == http.StatusOK {
				f(resp.String())
			}
		}

	}()
		
}

func Publish(topic string, data string) error {
	return errors.New("not implemented")
}