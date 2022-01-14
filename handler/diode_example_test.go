package handler_test

import (
	"fmt"
	"os"

	"github.com/covine/easylog"
)

func ExampleNewWriter() {
	_ = easylog.GetRootLogger()

	w := NewWriter(os.Stdout, 1000, 0, func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	})
	log := zerolog.New(w)
	log.Print("test")

	w.Close()

	// Output: {"level":"debug","message":"test"}
}
