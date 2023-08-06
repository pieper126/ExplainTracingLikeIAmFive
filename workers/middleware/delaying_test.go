package middleware_test

import (
	"context"
	"testing"
	"time"

	"ExplainTracingLikeIAmFive/workers"
	"ExplainTracingLikeIAmFive/workers/middleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelayedMiddlewareDelaysMessage(t *testing.T) {
	messageDelay := time.Millisecond * 100
	delayingMiddleware := middleware.NewDelaying("test", messageDelay)
	require.NotNil(t, delayingMiddleware)

	messageInterval := time.Millisecond * 1
	timeout := time.Millisecond * 10

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	input := messageOnInterval(ctx, messageInterval)

	outputFromMiddleware := delayingMiddleware.Run(ctx, input)

	delivered := make([]workers.MessageWithContext, 0)

	for message := range outputFromMiddleware {
		delivered = append(delivered, message)
	}

	assert.Len(t, delivered, 1)
}

func messageOnInterval(ctx context.Context, d time.Duration) chan workers.MessageWithContext {
	tick := time.Tick(d)

	return mapTimeStampToMessageWithContext(ctx, tick)
}

func mapTimeStampToMessageWithContext(ctx context.Context, in <-chan time.Time) chan workers.MessageWithContext {
	res := make(chan workers.MessageWithContext, 10)

	go func(input <-chan time.Time, out chan workers.MessageWithContext) {
		for {
			select {
			case <-ctx.Done():
				close(res)
				return
			case message := <-input:
				out <- workers.MessageWithContext{
					Message: message.Format(time.RFC3339),
					Ctx:     context.Background(),
				}
			}
		}
	}(in, res)

	return res
}
