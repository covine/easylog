package handler

import (
	"testing"

	"github.com/covine/easylog"
)

func TestStdoutHandlerWithLineFormatter(t *testing.T) {
	defer easylog.Flush()

	easylog.EnableCaller(easylog.INFO)
	defer easylog.DisableCaller(easylog.INFO)
	easylog.EnableCaller(easylog.ERROR)
	defer easylog.DisableCaller(easylog.ERROR)

	easylog.EnableStack(easylog.ERROR)
	defer easylog.DisableStack(easylog.ERROR)

	h := NewStdoutHandler(StdFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	easylog.Debug().Log()
	easylog.Debug().Logf("stdout")

	easylog.Info().Log()
	easylog.Info().Logf("stdout")

	easylog.Warn().Log()
	easylog.Warn().Logf("error")

	easylog.Error().Logf("error")

	logger := easylog.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(easylog.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(easylog.ERROR)

	logger.Error().Logf("error")
}

func TestStdoutHandlerWithJsonFormatter(t *testing.T) {
	defer easylog.Flush()

	easylog.EnableCaller(easylog.INFO)
	defer easylog.DisableCaller(easylog.INFO)
	easylog.EnableCaller(easylog.ERROR)
	defer easylog.DisableCaller(easylog.ERROR)

	easylog.EnableStack(easylog.ERROR)
	defer easylog.DisableStack(easylog.ERROR)

	h := NewStdoutHandler(JsonFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	easylog.Debug().Log()
	easylog.Debug().Logf("stdout")

	easylog.Info().Log()
	easylog.Info().Logf("stdout")

	easylog.Warn().Log()
	easylog.Warn().Logf("error")

	easylog.Error().Logf("error")

	logger := easylog.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(easylog.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(easylog.ERROR)

	logger.Error().Logf("error")
}

func TestStderrHandlerWithLineFormatter(t *testing.T) {
	defer easylog.Flush()

	easylog.EnableCaller(easylog.INFO)
	defer easylog.DisableCaller(easylog.INFO)
	easylog.EnableCaller(easylog.ERROR)
	defer easylog.DisableCaller(easylog.ERROR)

	easylog.EnableStack(easylog.ERROR)
	defer easylog.DisableStack(easylog.ERROR)

	h := NewStderrHandler(StdFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	easylog.Debug().Log()
	easylog.Debug().Logf("stdout")

	easylog.Info().Log()
	easylog.Info().Logf("stdout")

	easylog.Warn().Log()
	easylog.Warn().Logf("error")

	easylog.Error().Logf("error")

	logger := easylog.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(easylog.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(easylog.ERROR)

	logger.Error().Logf("error")
}

func TestStderrHandlerWithJsonFormatter(t *testing.T) {
	defer easylog.Flush()

	easylog.EnableCaller(easylog.INFO)
	defer easylog.DisableCaller(easylog.INFO)
	easylog.EnableCaller(easylog.ERROR)
	defer easylog.DisableCaller(easylog.ERROR)

	easylog.EnableStack(easylog.ERROR)
	defer easylog.DisableStack(easylog.ERROR)

	h := NewStderrHandler(JsonFormatter)

	easylog.AddHandler(h)
	defer easylog.RemoveHandler(h)

	easylog.Debug().Log()
	easylog.Debug().Logf("stdout")

	easylog.Info().Log()
	easylog.Info().Logf("stdout")

	easylog.Warn().Log()
	easylog.Warn().Logf("error")

	easylog.Error().Logf("error")

	logger := easylog.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(easylog.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(easylog.ERROR)

	logger.Error().Logf("error")
}
