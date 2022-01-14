package handler

import (
	"context"
	"time"

	"github.com/covine/easylog"
	"github.com/covine/easylog/diode"
)

type puller interface {
	diode.Diode

	Next() diode.GenericDataType
}

type RingBufferHandler struct {
	handler    easylog.Handler
	ringBuffer diode.Diode
	puller     puller
	cancel     context.CancelFunc
	done       chan struct{}
}

func RingBufferWrapper(
	h easylog.Handler, size int, interval time.Duration, alert diode.AlertFunc,
) easylog.Handler {
	ctx, cancel := context.WithCancel(context.Background())

	r := &RingBufferHandler{
		cancel: cancel,
		done:   make(chan struct{}),
	}

	d := diode.NewManyToOne(size, alert)
	if interval > 0 {
		r.puller = diode.NewPoller(
			d,
			diode.WithPollingInterval(interval),
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
	r.puller.Set(diode.GenericDataType(e))

	return true, nil
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

	// TODO r.handler.close()
	// writer?

	return nil
}

func (r *RingBufferHandler) pull() {
	defer close(r.done)

	for {
		d := r.puller.Next()
		if d == nil {
			return
		}

		_ = (*easylog.Event)(d)
		// TODO handler e
		// write?
		// another event pool here?
	}
}

// another pool here?
