package tracer

import (
	"context"
	"github.com/ozonmp/lic-license-api/internal/pkg/logger"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// NewTracer - returns new tracer.
func NewTracer(ctx context.Context, serviceName string, host string, port string) (io.Closer, error) {
	cfgTracer := &jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: host + port,
		},
	}
	tracer, closer, err := cfgTracer.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		logger.ErrorKV(ctx, "failed init jaeger", "err", err)

		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	logger.InfoKV(ctx, "traces started")

	return closer, nil
}

type span struct {
	base opentracing.Span
}

// StartSpanFromContext - стартует подчиненный спан.
func StartSpanFromContext(ctx context.Context, name string) *span {
	baseSp := opentracing.SpanFromContext(ctx)

	var sp opentracing.Span
	if baseSp == nil {
		sp = opentracing.StartSpan(name)
	} else {
		sp = opentracing.StartSpan(name, opentracing.FollowsFrom(baseSp.Context()))
	}

	return &span{
		base: sp,
	}
}

// Finish - заканчивает спан.
func (s *span) Finish() {
	s.base.Finish()
}
