package middleware_test

import (
	"ExplainTracingLikeIAmFive/workers"
	"ExplainTracingLikeIAmFive/workers/middleware"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntermittentlyBroken100Procent(t *testing.T) {
	timeout := time.Millisecond * 5
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	name := "test"
	percentageBroken := 100.0

	brokenMiddleware := middleware.NewIntermittentlyBroken(name, "100%Broken", percentageBroken, assert.AnError)
	require.NotNil(t, brokenMiddleware)

	in := make(chan workers.MessageWithContext, 2)

	in <- workers.MessageWithContext{
		Ctx: context.Background(),
	}
	in <- workers.MessageWithContext{
		Ctx: context.Background(),
	}

	output := brokenMiddleware.Run(ctx, in)

	delivered := make([]workers.MessageWithContext, 0)

	for message := range output {
		delivered = append(delivered, message)
	}

	assert.Len(t, delivered, 0, "with a 100 percentage failure none should pass")
	cancel()
}

func TestIntermittentlyBroken0Procent(t *testing.T) {
	timeout := time.Millisecond * 5
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	name := "test"
	percentageBroken := 0.0

	brokenMiddleware := middleware.NewIntermittentlyBroken(name, "100%Broken", percentageBroken, assert.AnError)
	require.NotNil(t, brokenMiddleware)

	in := make(chan workers.MessageWithContext, 2)

	in <- workers.MessageWithContext{
		Ctx: context.Background(),
	}
	in <- workers.MessageWithContext{
		Ctx: context.Background(),
	}

	output := brokenMiddleware.Run(ctx, in)

	delivered := make([]workers.MessageWithContext, 0)

	for message := range output {
		delivered = append(delivered, message)
	}

	assert.Len(t, delivered, 2, "with a 0 percentage failure should pass everything")
	cancel()
}

func TestIntermittentlyBroken50Percent(t *testing.T) {
	timeout := time.Millisecond * 5
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	name := "test"
	percentageBroken := 1.0

	brokenMiddleware := middleware.NewIntermittentlyBroken(name, "100%Broken", percentageBroken, assert.AnError)
	require.NotNil(t, brokenMiddleware)

	numberOfMessages := 100

	in := make(chan workers.MessageWithContext, numberOfMessages)

	for i := 0; i < numberOfMessages; i++ {
		in <- workers.MessageWithContext{
			Ctx: context.Background(),
		}
	}

	output := brokenMiddleware.Run(ctx, in)

	delivered := make([]workers.MessageWithContext, 0)

	for message := range output {
		delivered = append(delivered, message)
	}

	assert.Greater(t, len(delivered), 0)
	assert.True(t, len(delivered) < numberOfMessages)
	cancel()
}
