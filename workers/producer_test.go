package workers_test

import (
	"context"
	"testing"
	"time"

	"ExplainTracingLikeIAmFive/workers"

	"github.com/stretchr/testify/assert"
)

func TestAbleToCancel(t *testing.T) {
	producer := workers.NewProducer("producer-1")

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	in := producer.Run(ctx, time.Millisecond*40, func() string { return "0" })

	t1 := <-in
	t2 := <-in

	assert.Equal(t, "0", t1.Message)
	assert.Equal(t, "0", t2.Message)
	cancel()

	<-in

	assert.True(t, true)
}

type MessageWithTimeStamp struct {
	Message   string
	Timestamp time.Time
}

func TestIfTheIntervalIsRespective(t *testing.T) {
	producer := workers.NewProducer("producer-1")

	timeout := time.Millisecond * 10
	intervalPerMessage := time.Millisecond * 1

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	in := producer.Run(ctx, intervalPerMessage, func() string { return "0" })

	sent := make([]MessageWithTimeStamp, 0)

	for message := range in {
		sent = append(sent, MessageWithTimeStamp{
			Message:   message.Message,
			Timestamp: time.Now(),
		})
	}

	assert.Len(t, sent, int(timeout.Milliseconds())/int(intervalPerMessage.Milliseconds()))
}
