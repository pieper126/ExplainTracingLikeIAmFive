package middleware

import (
	"ExplainTracingLikeIAmFive/workers"
	"context"
	"math/rand"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type IntermittentlyBroken struct {
	name             string
	nameOperation    string
	withError        error
	percentageBroken float64
}

func NewIntermittentlyBroken(name string, nameOperation string, percentageBroken float64, withError error) *IntermittentlyBroken {
	return &IntermittentlyBroken{
		name:             name,
		percentageBroken: percentageBroken,
		nameOperation:    nameOperation,
		withError:        withError,
	}
}

func (d *IntermittentlyBroken) Run(ctx context.Context, input <-chan workers.MessageWithContext) chan workers.MessageWithContext {
	res := make(chan workers.MessageWithContext, 10)

	go func(in <-chan workers.MessageWithContext, out chan workers.MessageWithContext) {
		innerTr := otel.Tracer(d.name)
		for {
			select {
			case <-ctx.Done():
				close(res)
				return
			case incoming := <-in:
				_, span := innerTr.Start(incoming.Ctx, d.nameOperation)
				if d.shouldSent() {
					out <- incoming
				} else {
					span.RecordError(d.withError)
					span.SetStatus(codes.Error, d.withError.Error())
				}
				span.End()
			}
		}
	}(input, res)

	return res
}

func (d *IntermittentlyBroken) shouldSent() bool {
	if d.percentageBroken == 0 {
		return true
	}

	normalised := (d.percentageBroken / 100)
	f := rand.Float64()

	return f > normalised
}
