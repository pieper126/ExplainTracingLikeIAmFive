package workers_test

import (
	"ExplainTracingLikeIAmFive/workers"
	"ExplainTracingLikeIAmFive/workers/middleware"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoMiddleWare(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	listOfMiddleware := []workers.Middleware{}
	pipeline := workers.NewPipeline(listOfMiddleware)
	require.NotNil(t, pipeline)

	in := make(chan workers.MessageWithContext, 10)

	out := pipeline.Run(ctx, in)

	testValue := workers.MessageWithContext{
		Message: "test",
		Ctx:     context.Background(),
	}

	in <- testValue

	res := <-out

	assert.Equal(t, testValue, res)
}

func TestOneMiddleWare(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	delaying := middleware.NewDelaying("test delaying", time.Millisecond*10)

	listOfMiddleware := []workers.Middleware{delaying}
	pipeline := workers.NewPipeline(listOfMiddleware)
	require.NotNil(t, pipeline)

	in := make(chan workers.MessageWithContext, 10)

	out := pipeline.Run(ctx, in)

	testValue := workers.MessageWithContext{
		Message: "test",
		Ctx:     context.Background(),
	}

	in <- testValue

	res := <-out

	ctx.Done()

	assert.Equal(t, testValue.Message, res.Message)
}
