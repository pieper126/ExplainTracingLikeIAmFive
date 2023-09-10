package middleware

import (
	"context"
	"time"

	"ExplainTracingLikeIAmFive/workers"

	"go.opentelemetry.io/otel"
)

type Delaying struct {
	name  string
	delay time.Duration
}

func NewDelaying(name string, delay time.Duration) *Delaying {
	return &Delaying{
		name:  name,
		delay: delay,
	}
}

func (d *Delaying) Run(ctx context.Context, input <-chan workers.MessageWithContext) chan workers.MessageWithContext {
	res := make(chan workers.MessageWithContext, 10)

	go func(in <-chan workers.MessageWithContext, out chan workers.MessageWithContext) {
		innerTr := otel.Tracer("middle-1")
		for {
			select {
			case <-ctx.Done():
				close(res)
				return
			case incoming := <-in:
				ctx, span := innerTr.Start(incoming.Ctx, "processing")
				time.Sleep(d.delay)
				out <- workers.MessageWithContext{
					Message: incoming.Message,
					Ctx:     ctx,
				}
				span.End()
			}
		}
	}(input, res)

	return res
}
