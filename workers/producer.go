package workers

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
)

type MessageWithContext struct {
	Message string
	Ctx     context.Context
}

type Producer struct {
	name string
}

func NewProducer(name string) *Producer {
	return &Producer{
		name: name,
	}
}

func (p *Producer) Run(ctx context.Context, onInterval time.Duration, messageFn func() string) chan MessageWithContext {
	res := make(chan MessageWithContext, 10)

	go func(localCtx context.Context, out chan MessageWithContext) {
		innerTr := otel.Tracer(p.name)
		for {
			select {
			case <-localCtx.Done():
				close(res)
				return
			case <-time.After(onInterval):
				ctx, bla := innerTr.Start(context.Background(), "producing-of-message")
				out <- MessageWithContext{
					Message: messageFn(),
					Ctx:     ctx,
				}
				bla.End()
			}
		}
	}(ctx, res)

	return res
}
