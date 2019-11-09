package main

import (
	"context"
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
	defaultDialTimeout   = time.Second
	defaultDialAddress   = "127.0.0.1:2379"
	defaultReadTimeout   = 250 * time.Millisecond
	defaultListenAddress = "127.0.0.1:9000"
	requestTimeout       = 2 * time.Second
)

func buildCmdTree(logger *log.Logger, _ io.Writer, rootCmd *cobra.Command) {
	rootCmd.RunE = func(_ *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		return run(ctx, logger)
	}
}

func run(_ context.Context, logger *log.Logger) error {
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

	logger.Printf("Starting server at address:%s", defaultListenAddress)
	err = errors.Wrap((&osc2.Server{
		Addr: defaultListenAddress,
		Dispatcher: osc.Dispatcher{
			KeyValuePutter: etcd.Client{KV: clientv3.NewKV(cli)},
			HandleError:    loggingErrorHandler(logger),
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

func client(dialTimeout time.Duration, addr string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{addr},
	})
}
