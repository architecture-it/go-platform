package mq

import (
	"context"
	"gopkg.in/resty.v1"
	"fmt"
	"errors"
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
		
		for {
			select {
			case <- ctx.Done():
				return
			}
			resp,err := resty.R().Get(fmt.Sprintf("%s/queues/%s",q.api,q.name))
			if err == nil {
				f(resp.String())
			}
		}

	}()
		
}

func Publish(topic string, data string) error {
	return errors.New("not implemented")
}