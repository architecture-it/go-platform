package AMQStream

import (
	"sync"
)

var lock = &sync.Mutex{}
var singleInstance *config

func getInstance() *config {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &config{}
			singleInstance.consumers = []ConsumerOptions{}
		}
	}

	return singleInstance
}
