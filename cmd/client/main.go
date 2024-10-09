package main

import (
	"TinyKVStore/internal/composer"
	"TinyKVStore/internal/helpers/syncutil"
	"log"
	"syscall"
)

func main() {
	app, err := composer.ComposeNewClientApplication()
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalln(err)
		}
	}()

	syncutil.WaitForSignal(syscall.SIGINT)

	if err := app.Shutdown(); err != nil {
		log.Fatalln(err)
	}
}
