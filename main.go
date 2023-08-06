// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command jaeger is an example program that creates spans
// and uploads to Jaeger.
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"ExplainTracingLikeIAmFive/workers"
	"ExplainTracingLikeIAmFive/workers/middleware"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	service        = "trace-demo-1"
	environment    = "production"
	id             = 1
	defaultTimeOut = time.Second * 5
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

func main() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	// tr := otel.Tracer("test")

	// spanCtx, span := tr.Start(ctx, "wrapping-span")

	// PingEverySecond(spanCtx)

	// span.End()
	// TestTracing()
	UsingNewThigns()
}

func UsingNewThigns() {
	groupContext, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	producer1 := workers.NewProducer("producer")
	delay1 := middleware.NewDelaying("delay 1", time.Microsecond*10)
	mightBeBroken := middleware.NewIntermittentlyBroken("broken 1", "query", 50, fmt.Errorf("random error"))
	endpoint := workers.NewEndpoint("endpoint")

	in := producer1.Run(groupContext, time.Millisecond*10, func() string { return "ping" })
	outDelay1 := delay1.Run(groupContext, in)
	outDelay2 := mightBeBroken.Run(groupContext, outDelay1)
	endpoint.Run(
		groupContext,
		func(mwc workers.MessageWithContext) { fmt.Println("pong") },
		outDelay2,
	)

	<-groupContext.Done()
}

func PingEverySecond(ctx context.Context) {
	timeout := time.Tick(time.Second * 10)
	second := time.Tick(time.Second * 1)

	var wg sync.WaitGroup

	wg.Add(1)

	go func(routineCtx context.Context) {
		innerTr := otel.Tracer("pinggingRoutine")

		_, innerSpan := innerTr.Start(routineCtx, "start inner pinging", trace.WithNewRoot())

		for {
			select {
			case <-timeout:
				innerSpan.End()
				wg.Done()
				return
			case <-second:
				innerSpan.AddEvent("ping", trace.WithTimestamp(time.Now()))
			}
		}
	}(ctx)

	wg.Wait()
}

type MessageWithContext struct {
	Message string
	ctx     context.Context
}

// blocking
func TestTracing() {
	var wg sync.WaitGroup

	// producer 1
	wg.Add(1)
	ticker := time.Tick(time.Second * 1)
	prOut := make(chan MessageWithContext)
	go func(in <-chan time.Time, out chan MessageWithContext) {
		innerTr := otel.Tracer("producer-1")
		for {
			select {
			case <-time.After(defaultTimeOut):
				wg.Done()
				return
			case incoming := <-in:
				fmt.Println(incoming)
				ctx, bla := innerTr.Start(context.Background(), "producing-of-message")
				out <- MessageWithContext{
					Message: incoming.Format(time.RFC3339),
					ctx:     ctx,
				}
				bla.End()
			}
		}
	}(ticker, prOut)

	// middle
	wg.Add(1)
	middleOut := make(chan MessageWithContext, 10)
	go func(in chan MessageWithContext, out chan MessageWithContext) {
		innerTr := otel.Tracer("middle-1")
		for {
			select {
			case <-time.After(defaultTimeOut):
				wg.Done()
				return
			case incoming := <-in:
				fmt.Println("middle")
				ctx, bla := innerTr.Start(incoming.ctx, "processing")
				time.Sleep(time.Second * 1)
				out <- MessageWithContext{
					Message: incoming.Message,
					ctx:     ctx,
				}
				bla.End()
			}
		}
	}(prOut, middleOut)

	// end
	wg.Add(1)
	go func(in chan MessageWithContext) {
		innerTr := otel.Tracer("end-1")
		for {
			select {
			case <-time.After(defaultTimeOut):
				fmt.Println("end stopping")
				wg.Done()
				return
			case incoming := <-in:
				_, span := innerTr.Start(incoming.ctx, "end")
				fmt.Printf("end for message: %s \n", incoming.Message)
				span.End()
			}
		}
	}(middleOut)

	wg.Wait()
}
