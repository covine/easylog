package handler

import (
	"context"
	"sync"
	"time"

	"github.com/covine/easylog"
	"github.com/covine/easylog/diode"
)

type diodePuller interface {
	diode.Diode

	Next() diode.GenericDataType
}

type RingBufferHandler struct {
	handler  easylog.Handler
	puller   diodePuller
	interval time.Duration
	cancel   context.CancelFunc
	done     chan struct{}
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
	// p is pooled in zerolog so we can't hold it passed this call, hence the
	// copy.
	//p = append(bufPool.Get().([]byte), p...)
	//dw.d.Set(diode.GenericDataType(&p))
	//return len(p), nil

	return true, nil
}

func (r *RingBufferHandler) Flush() error {
	return nil
}

func (r *RingBufferHandler) Close() error {
	r.cancel()
	<-r.done

	//if w, ok := dw.w.(io.Closer); ok {
	//return w.Close()
	//}

	return nil
}

func (r *RingBufferHandler) pull() {
	defer close(r.done)

	for {
		d := r.puller.Next()
		if d == nil {
			return
		}

		//p := *(*[]byte)(d)
		//dw.w.Write(p)

		// Proper usage of a sync.Pool requires each entry to have approximately
		// the same memory cost. To obtain this property when the stored type
		// contains a variably-sized buffer, we add a hard limit on the maximum buffer
		// to place back in the pool.
		//
		// See https://golang.org/issue/23199
		//const maxSize = 1 << 16 // 64KiB
		//if cap(p) <= maxSize {
		//bufPool.Put(p[:0])
		//}
	}
}

var bufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 500)
	},
}
