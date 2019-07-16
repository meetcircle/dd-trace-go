package mgo

import (
	"github.com/globalsign/mgo"
	"github.com/meetcircle/dd-trace-go/ddtrace/tracer"
)

// Pipe is an mgo.Pipe instance along with the data necessary for tracing.
type Pipe struct {
	*mgo.Pipe
	cfg  *mongoConfig
	tags map[string]string
}

// Iter invokes and traces Pipe.Iter
func (p *Pipe) Iter() *Iter {
	span := newChildSpanFromContext(p.cfg, p.tags)
	iter := p.Pipe.Iter()
	span.Finish()
	return &Iter{
		Iter: iter,
		cfg:  p.cfg,
		tags: p.tags,
	}
}

// All invokes and traces Pipe.All
func (p *Pipe) All(result interface{}) error {
	return p.Iter().All(result)
}

// One invokes and traces Pipe.One
func (p *Pipe) One(result interface{}) (err error) {
	span := newChildSpanFromContext(p.cfg, p.tags)
	defer span.Finish(tracer.WithError(err))
	err = p.Pipe.One(result)
	return
}

// Explain invokes and traces Pipe.Explain
func (p *Pipe) Explain(result interface{}) (err error) {
	span := newChildSpanFromContext(p.cfg, p.tags)
	defer span.Finish(tracer.WithError(err))
	err = p.Pipe.Explain(result)
	return
}
