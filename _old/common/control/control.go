package control

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pavlo67/common/common/logger"
)

var signalChan = make(chan os.Signal, 1000)

func Init(l logger.Operator) error {
	if l == nil {
		return fmt.Errorf("on control.Init: no logger")
	}
	signal.Notify(signalChan, syscall.SIGPIPE)

	go processSignal(l)

	return nil
}

func processSignal(l logger.Operator) {
	for {
		sig := <-signalChan

		if sig == syscall.SIGPIPE {
			signal.Reset(sig)
			l.Warnf("got & ignored signal %s", sig)
		} else {
			l.Warnf("got signal %s", sig)
		}
	}
}
