package testserver

import (
	"context"

	"github.com/pkg/errors"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/worker"
)

// Env encapsulates DevServer related propoerties and keeps state.
// usage:
//   - use New to create env
//   - env.Start to prepare the environment and start the DevServer
//   - run workflow code as one would against the actual temporal environment
//   - env.Stop to cleanup and stop DeverServer
type Env struct {
	ServerOpts testsuite.DevServerOptions
	WorkerOpts worker.Options
	TaskQ      string

	server *testsuite.DevServer
	worker worker.Worker
}

// New create a Env with specified properties
//
//	hostPort   - "host:port" form. Default "" will let use a random <port> in your local env
//	taskQ      - taskQ for the workflows
func New(hostPort, taskQ string) *Env {
	e := &Env{
		ServerOpts: testsuite.DevServerOptions{ClientOptions: &client.Options{HostPort: hostPort}},
		TaskQ:      taskQ,
	}

	return e
}

func (e *Env) startServer() error {
	var err error

	ctx := context.Background()
	e.server, err = testsuite.StartDevServer(ctx, e.ServerOpts)
	if err != nil {
		return err
	}

	return nil
}

func (e *Env) startWorker() error {
	ch := make(chan error)
	go func() {
		c := e.server.Client()

		e.worker = worker.New(c, e.TaskQ, worker.Options{})

		ch <- nil

		_ = e.worker.Run(worker.InterruptCh())
	}()

	return <-ch
}

// Start create and start the server and worker
func (e *Env) Start() error {
	err := e.startServer()
	if err != nil {
		return err
	}

	return e.startWorker()
}

// Client get the client created as part of the DevServer's startup
func (e *Env) Client() (client.Client, error) {
	if e.worker == nil {
		return nil, errors.New("worker not started")
	}

	return e.server.Client(), nil
}

// Worker get the worker created as part of the DevServer's startup
func (e *Env) Worker() (worker.Worker, error) {
	if e.server == nil {
		return nil, errors.New("server not started")
	}

	return e.worker, nil
}

// Stop stops the server if one was created
func (e *Env) Stop() error {
	if e.server != nil {
		return e.server.Stop()
	}

	return nil
}
