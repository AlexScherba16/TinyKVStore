package main

import (
	"TinyKVStore/internal/composer"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	app, err := composer.ComposeNewClientApplication()
	if err != nil {
		panic(err)
	}

	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	<-signalChan
	if err := app.Shutdown(); err != nil {
		panic(err)
	}
}
