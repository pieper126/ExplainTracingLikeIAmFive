package middleware

import (
	"context"

	"ExplainTracingLikeIAmFive/workers"

	"go.opentelemetry.io/otel"
)

type MapperFn = func(context.Context, workers.MessageWithContext) workers.MessageWithContext

type Mapper struct {
	name string
	fn   MapperFn
}

func NewMapper(name string, fn MapperFn) *Mapper {
	return &Mapper{
		name: name,
		fn:   fn,
	}
}

func (d *Mapper) Run(ctx context.Context, input <-chan workers.MessageWithContext) chan workers.MessageWithContext {
	res := make(chan workers.MessageWithContext, 10)

	go func(in <-chan workers.MessageWithContext, out chan workers.MessageWithContext) {
		innerTr := otel.Tracer("mapper-1")
		for {
			select {
			case <-ctx.Done():
				close(res)
				return
			case incoming := <-in:
				ctx, span := innerTr.Start(incoming.Ctx, "processing")
				out <- d.fn(ctx, incoming)
				span.End()
			}
		}
	}(input, res)

	return res
}
