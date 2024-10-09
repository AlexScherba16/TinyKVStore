package client

import (
	"TinyKVStore/internal/application"
	"TinyKVStore/internal/logger"
	"TinyKVStore/internal/network"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

const (
	clientShutdownDuration = time.Duration(1 * time.Minute)
)

// kvClient private structure defines a client application.
type kvClient struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        *sync.WaitGroup
	logger    logger.ILogger
	tcpClient *network.TCPClient
}

// NewClientApplication initializes and returns a new client application instance that implements the IApplication interface.
// Returns the initialized client application and any error encountered during setup.
func NewClientApplication(iLogger logger.ILogger, tcpClient *network.TCPClient) (application.IApplication, error) {
	if iLogger == nil {
		return nil, errors.New("invalid client logger argument")
	}

	if tcpClient == nil {
		return nil, errors.New("invalid tcp client argument")
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &kvClient{
		ctx:       ctx,
		cancel:    cancel,
		wg:        &sync.WaitGroup{},
		logger:    iLogger,
		tcpClient: tcpClient,
	}, nil
}

func (client *kvClient) run() {
	client.wg.Add(1)
	defer client.wg.Done()

	reader := bufio.NewReader(os.Stdin)
	clientPrompt := color.HiGreenString("[KV client] > ")

	for {
		fmt.Print(clientPrompt)

		request, err := reader.ReadString('\n')
		if err != nil {

			if err == io.EOF {
				client.logger.Debug("exit client application loop")
				return
			}
			client.logger.Error(err.Error())
			continue
		}

		err = client.tcpClient.Open()
		if err != nil {
			client.logger.Error(err.Error())
			continue
		}

		defer func() {
			client.tcpClient.Close()
		}()

		err = client.tcpClient.Write([]byte(request))
		if err != nil {
			client.logger.Error(err.Error())
			continue
		}

		response, err := client.tcpClient.Read()
		if err != nil {
			client.logger.Error(err.Error())
			continue
		}

		fmt.Println(response)
	}
	//if errors.Is(err, syscall.EPIPE) {
	//	logger.Fatal("connection was closed", zap.Error(err))
	//} else if err != nil {
	//	logger.Error("failed to read query", zap.Error(err))
	//}
}

// Run starts the client application.
// Returns error if something went wrong during initialization.
func (client *kvClient) Run() error {
	client.logger.Info("Press (Ctrl+C) to shutdown application")

	go client.run()

	<-client.ctx.Done()
	return nil
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
