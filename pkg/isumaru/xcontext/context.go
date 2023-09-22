package xcontext

import (
	"context"
	"time"
)

type Now struct{}

type keyConstraint interface {
	Now
}

type valueConstraint interface {
	time.Time
}

type key[T keyConstraint] struct{}

func WithValue[k keyConstraint, v valueConstraint](ctx context.Context, val v) context.Context {
	return context.WithValue(ctx, key[k]{}, val)
}

func Value[k keyConstraint, v valueConstraint](ctx context.Context) v {
	return ctx.Value(key[k]{}).(v)
}
