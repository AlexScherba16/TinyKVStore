package syncutil

import (
	"os"
	"os/signal"
)

type Emit func()
type Await func()

func WaitForSignal(sig os.Signal) {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)

	signal.Notify(signalChan, sig)
	<-signalChan
}

func NewEmitterAwaiter() (Emit, Await) {
	ch := make(chan struct{})
	emit := func() {
		close(ch)
	}

	await := func() {
		<-ch
	}
	return emit, await
}
