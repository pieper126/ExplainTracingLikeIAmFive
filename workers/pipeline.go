package workers

import "context"

type Middleware interface {
	Run(ctx context.Context, input <-chan MessageWithContext) chan MessageWithContext
}

type Pipeline struct {
	middleWares []Middleware
}

func NewPipeline(middleWares []Middleware) *Pipeline {
	return &Pipeline{
		middleWares: middleWares,
	}
}

func (p *Pipeline) Run(ctx context.Context, input <-chan MessageWithContext) <-chan MessageWithContext {
	if len(p.middleWares) == 0 {
		return input
	}

	res := p.initiatePipeline(ctx, input)

	return res
}

func (p *Pipeline) initiatePipeline(ctx context.Context, input <-chan MessageWithContext) <-chan MessageWithContext {
	middle := input

	for _, v := range p.middleWares {
		middle = v.Run(ctx, middle)
	}

	return middle
}
