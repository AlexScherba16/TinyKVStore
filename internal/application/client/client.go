package client

import (
	"TinyKVStore/internal/application"
	"TinyKVStore/internal/logger"
	"context"
	"errors"
	"sync"
	"time"
)

const (
	clientShutdownDuration = time.Duration(1 * time.Minute)
)

// kvClient private structure defines a client application.
type kvClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	logger logger.ILogger
}

// NewClientApplication initializes and returns a new client application instance that implements the IApplication interface.
// Returns the initialized client application and any error encountered during setup.
func NewClientApplication(iLogger logger.ILogger) (application.IApplication, error) {
	if iLogger == nil {
		return nil, errors.New("invalid client logger argument")
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &kvClient{
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
		logger: iLogger,
	}, nil
}

// Run starts the client application.
// Returns error if something went wrong during initialization.
func (client *kvClient) Run() error {
	client.logger.Info("Press (Ctrl+C) to shutdown application")

	client.wg.Add(1)
	defer client.wg.Done()

	// TODO: implement clients logic here

	// Dummy code
	for {
		select {
		case <-client.ctx.Done():
			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

// Shutdown gracefully shuts down the client application.
// Returns error if something went wrong during shutdown.
func (client *kvClient) Shutdown() error {
	applicationDoneChan := make(chan struct{})
	go func() {
		client.wg.Wait()
		close(applicationDoneChan)
	}()
	client.cancel()

	select {
	case <-time.After(clientShutdownDuration):
		return errors.New("shutdown timeout")
	case <-applicationDoneChan:
		return nil
	}
}
