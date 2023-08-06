package workers

import (
	"context"

	"go.opentelemetry.io/otel"
)

type Endpoint struct {
	name string
}

func NewEndpoint(name string) *Endpoint {
	return &Endpoint{
		name: name,
	}
}

func (c *Endpoint) Run(ctx context.Context, onMessage func(MessageWithContext), in chan MessageWithContext) {
	go func(in chan MessageWithContext) {
		innerTr := otel.Tracer(c.name)
		for {
			select {
			case <-ctx.Done():
				return
			case incoming := <-in:
				_, span := innerTr.Start(incoming.Ctx, "end")
				onMessage(incoming)
				span.End()
			}
		}
	}(in)
}
