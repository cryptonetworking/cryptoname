package ctx

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"time"
	"unsafe"
)

type CancelFunc func(err error)

type ctx struct {
	ctx    context.Context
	cancel context.CancelFunc
	err    atomic.UnsafePointer
}

func (ctx *ctx) String() string {
	return ctx.ctx.(fmt.Stringer).String()
}

func (ctx *ctx) Deadline() (deadline time.Time, ok bool) {
	return ctx.ctx.Deadline()
}

func (ctx *ctx) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

func (ctx *ctx) Err() error {
	err := ctx.err.Load()
	if err == nil {
		return nil
	}
	return *((*error)(err))
}

func (ctx *ctx) Value(key any) any {
	return ctx.ctx.Value(key)
}

func (ctx *ctx) Cancel(err error) {
	if ctx.err.CAS(nil, unsafe.Pointer(&err)) {
		ctx.cancel()
	}
}

func WithCancel(parent context.Context) (context.Context, CancelFunc) {
	c, cancel := context.WithCancel(parent)
	ctx := &ctx{ctx: c, cancel: cancel}
	return ctx, func(err error) {
		if err == nil {
			err = context.Canceled
		}
		ctx.Cancel(err)
	}
}
