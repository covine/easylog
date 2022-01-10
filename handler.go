package easylog

// Handler mockery --name=Handler --inpackage --case=underscore --filename=handler_mock.go --structname MockHandler
type Handler interface {
	// Handle consumes the Event, determine whether Event should be consumed by the next Handler and parent's Handler.
	// It returns a decision for the Logger to determine whether to continue handling the Event or not.
	// And returns an error which will be handled by the Logger's errHandler.
	// For example,
	// (true, <error>) means the Event will be consumed by the next Handler.
	// (false, <error>) means the Event will not be sequentially handled.
	Handle(*Event) (bool, error)
	// Flush flushes buffered data (if any).
	// Handler could have an internal flush() mechanism.
	// However, in some cases, only Logger knows when to flush().
	Flush() error
	// Close releases resources used by Handler. If any buffered data, should flush() before close().
	Close() error
}

// ErrorHandler mockery --name=ErrorHandler --inpackage --case=underscore --filename=error_handler_mock.go --structname MockErrorHandler
type ErrorHandler interface {
	// Handle consumes the error
	Handle(error) error
	// Flush flushes buffered data (if any).
	// ErrorHandler could have an internal flush() mechanism.
	// However, in some cases, only Logger knows when to flush().
	Flush() error
	// Close releases resources used by ErrorHandler. If any buffered data, should flush() before close().
	Close() error
}

type nopHandler struct {
}

func NewNopHandler() *nopHandler {
	return &nopHandler{}
}

func (n *nopHandler) Handle(_ *Event) (bool, error) {
	return true, nil
}

func (n *nopHandler) Flush() error {
	return nil
}

func (n *nopHandler) Close() error {
	return nil
}

type nopErrorHandler struct {
}

func (n *nopErrorHandler) Handle(_ error) error {
	return nil
}

func NewNopErrorHandler() *nopErrorHandler {
	return &nopErrorHandler{}
}

func (n *nopErrorHandler) Flush() error {
	return nil
}

func (n *nopErrorHandler) Close() error {
	return nil
}
