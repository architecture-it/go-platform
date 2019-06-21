package mq

import (
	"context"
	"testing"
	"time"
)

func TestBatchGet(t *testing.T) {

	q := GetQueue(ReadConfigFromEnv())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	i := 0
	q.Listen(ctx, func(msg string) {
		i++
	})
	<-ctx.Done()
	t.Log(i)
}

func BenchmarkBatchGet(b *testing.B) {

	q := GetQueue(ReadConfigFromEnv())
	ctx, cancel := context.WithCancel(context.Background())
	i := 0
	q.Listen(ctx, func(msg string) {
		i++
		if i > b.N {
			cancel()
		}
	})
	<-ctx.Done()
	b.Log(i)
}
