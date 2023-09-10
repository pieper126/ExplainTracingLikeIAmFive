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

func TestMapper(t *testing.T) {
	timeout := time.Millisecond * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	expectedValue := workers.MessageWithContext{
		Ctx:     ctx,
		Message: "pong",
	}

	delayingMiddleware := middleware.NewMapper(
		"test",
		func(ctx context.Context, mwc workers.MessageWithContext) workers.MessageWithContext {
			return expectedValue
		},
	)

	require.NotNil(t, delayingMiddleware)

	messageInterval := time.Millisecond * 1

	input := messageOnInterval(ctx, messageInterval)

	outputFromMiddleware := delayingMiddleware.Run(ctx, input)

	delivered := make([]workers.MessageWithContext, 0)

	for message := range outputFromMiddleware {
		delivered = append(delivered, message)
		break
	}

	assert.True(t, len(delivered) > 0)

	assert.Equal(t, expectedValue.Message, delivered[0].Message)
}
