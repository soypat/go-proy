package trace

import (
	"fmt"
	"io"

)

// Tracer is the interface that describes an object capable of
// tracing events throughout code.
type tracer struct {
	out io.Writer
}

type Tracer interface {
	Trace(...interface{})
}



func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type nilTracer struct{}
func (t *nilTracer) Trace(a ...interface{}) {}
// Off creates a Tracer that will ignore calls to Trace.
func Off() Tracer {
	return &nilTracer{}
}