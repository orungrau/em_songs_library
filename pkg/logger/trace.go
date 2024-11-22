package logger

import (
	"github.com/rs/zerolog"
	"time"
)

type TracingHook struct{}

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()

	if duration, ok := ctx.Value("duration").(time.Time); ok {
		e.Dur("duration", time.Since(duration))
	}

	if spanId, ok := ctx.Value("span-id").(string); ok && spanId != "" {
		e.Str("span-id", spanId)
	}
}
