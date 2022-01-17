package handler

import (
	"context"
	"time"

	"github.com/covine/easylog"
	"github.com/covine/easylog/diode"
)

type Puller interface {
	diode.Diode

	Next() diode.GenericDataType
}

type Next func(*easylog.Event) (bool, error)
type Handle func(*easylog.Event) error

// RingBufferHandler is a Ring Buffer Handler
type RingBufferHandler struct {
	next   Next
	handle Handle
	diode  diode.Diode
	puller Puller
	cancel context.CancelFunc
	done   chan struct{}
}

func RingBufferWrap(
	next Next, handle Handle, size int, alert diode.AlertFunc, pullInterval time.Duration,
) easylog.Handler {
	ctx, cancel := context.WithCancel(context.Background())

	r := &RingBufferHandler{
		next:   next,
		handle: handle,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	d := diode.NewManyToOne(size, alert)

	if pullInterval > 0 {
		r.puller = diode.NewPoller(
			d,
			diode.WithPollingInterval(pullInterval),
			diode.WithPollingContext(ctx),
		)
	} else {
		r.puller = diode.NewWaiter(
			d,
			diode.WithWaiterContext(ctx),
		)
	}

	go r.pull()

	return r
}

func (r *RingBufferHandler) Handle(e *easylog.Event) (bool, error) {
	r.puller.Set(diode.GenericDataType(e.Clone()))
	return r.next(e)
}

func (r *RingBufferHandler) Flush() error {
	for {
		data, ok := r.ringBuffer.TryNext()
		if !ok {
			return nil
		}

		e := (*easylog.Event)(data)
		// TODO handle event
		// TODO putEvent
		println(e)
	}
}

func (r *RingBufferHandler) Close() error {
	r.cancel()
	<-r.done

	return nil
}

func (r *RingBufferHandler) pull() {
	defer close(r.done)

	for {
		d := r.puller.Next()
		if d == nil {
			return
		}

		e := (*easylog.Event)(d)
		err := r.handle(e)
		if err != nil {
			// TODO errHandler
			e.Put()
		}
		e.Put()
	}
}
