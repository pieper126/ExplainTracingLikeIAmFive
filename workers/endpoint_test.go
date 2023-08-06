package workers_test

import (
	"context"
	"testing"
	"time"

	"ExplainTracingLikeIAmFive/workers"

	"github.com/stretchr/testify/assert"
)

func TestEndpointAbleToCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	endpoint := workers.NewEndpoint("test")

	producer := workers.NewProducer("test")
	in := producer.Run(ctx, time.Millisecond, func() string { return "test" })

	endpoint.Run(ctx, func(message workers.MessageWithContext) {}, in)

	cancel()

	assert.True(t, true)
}

func TestEndpointRunsOnMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	endpoint := workers.NewEndpoint("test")

	producer := workers.NewProducer("test")
	in := producer.Run(ctx, time.Millisecond, func() string { return "test" })

	collector := make([]workers.MessageWithContext, 0)

	endpoint.Run(ctx, func(message workers.MessageWithContext) {
		collector = append(collector, message)
	}, in)

	cancel()

	assert.GreaterOrEqual(t, 1, len(collector))
}
