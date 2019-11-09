package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/glynternet/etcd-osc-gateway/pkg/etcd"
	"github.com/glynternet/etcd-osc-gateway/pkg/osc"
	osc2 "github.com/glynternet/go-osc/osc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	defaultDialTimeout = time.Second
	defaultDialAddress = "127.0.0.1:2379"
	defaultReadTimeout = 250 * time.Millisecond
	requestTimeout     = 2 * time.Second
)

func buildCmdTree(logger *log.Logger, _ io.Writer, rootCmd *cobra.Command) {
	var listenHost string
	var listenPort uint

	rootCmd.RunE = func(_ *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		return run(ctx, logger, listenHost, listenPort)
	}

	rootCmd.Flags().StringVar(&listenHost, "listen-host", "127.0.0.1", "host address to listen on")
	rootCmd.Flags().UintVar(&listenPort, "listen-port", 9000, "host post to listen on")
}

func run(_ context.Context, logger *log.Logger, listenHost string, listenPort uint) error {
	cli, err := client(defaultDialTimeout, defaultDialAddress)
	if err != nil {
		return errors.Wrap(err, "creating client")
	}
	logger.Printf("Client created at address:%s", defaultDialAddress)
	defer func() {
		cErr := errors.Wrap(cli.Close(), "closing client")
		if cErr == nil {
			return
		}
		if err == nil {
			err = cErr
			return
		}
		logger.Println(cErr)
	}()

	// TODO: use Address type here
	listenAddress := fmt.Sprintf("%s:%d", listenHost, listenPort)
	logger.Printf("Starting server at address:%s", listenAddress)
	err = errors.Wrap((&osc2.Server{
		Addr: listenAddress,
		Dispatcher: osc.Dispatcher{
			KeyValuePutter: etcd.Client{KV: clientv3.NewKV(cli)},
			HandleError:    loggingErrorHandler(logger),
			HandleSuccess:  loggingSuccessHandler(logger),
		},
		ReadTimeout: defaultReadTimeout,
	}).ListenAndServe(), "serving")

	return err
}

func loggingErrorHandler(logger *log.Logger) func(error) {
	return func(err error) {
		logger.Println(err)
	}
}

func loggingSuccessHandler(logger *log.Logger) func(osc2.Message) {
	return func(msg osc2.Message) {
		logger.Println(msg)
	}
}

func client(dialTimeout time.Duration, addr string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{addr},
	})
}
