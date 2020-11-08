package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/glynternet/etcd-osc-gateway/pkg/etcd"
	"github.com/glynternet/etcd-osc-gateway/pkg/osc"
	osc2 "github.com/glynternet/go-osc/osc"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.etcd.io/etcd/client"
)

const (
	defaultDialTimeout = time.Second
	defaultDialAddress = "127.0.0.1:2379"
	defaultReadTimeout = 250 * time.Millisecond
	requestTimeout     = 2 * time.Second
)

func buildCmdTree(logger log.Logger, _ io.Writer, rootCmd *cobra.Command) {
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

func run(_ context.Context, logger log.Logger, listenHost string, listenPort uint) error {
	cli, err := etcdClient(defaultDialTimeout, defaultDialAddress)
	if err != nil {
		return errors.Wrap(err, "creating client")
	}
	if err := logger.Log(log.Message("Client created at address:%s"),
		log.KV{
			K: "dialAddress",
			V: defaultDialAddress,
		},
	); err != nil {
		return errors.Wrap(err, "logging during startup")
	}

	// TODO: use Address type here
	listenAddress := fmt.Sprintf("%s:%d", listenHost, listenPort)
	_ = logger.Log(log.Message("Starting server"), log.KV{
		K: "address",
		V: listenAddress,
	})

	err = errors.Wrap((&osc2.Server{
		Addr: listenAddress,
		Dispatcher: osc.Dispatcher{
			KeyValuePutter: etcd.Client{KeysAPI: client.NewKeysAPI(cli)},
			HandleError:    dispatchErrorLogger(logger),
			HandleSuccess:  loggingSuccessHandler(logger),
		},
		ReadTimeout: defaultReadTimeout,
	}).ListenAndServe(), "serving")

	return err
}

func dispatchErrorLogger(logger log.Logger) func(error) {
	return func(err error) {
		_ = logger.Log(log.Message("error dispatching message"), log.Error(err))
	}
}

func loggingSuccessHandler(logger log.Logger) func(osc2.Message) {
	return func(msg osc2.Message) {
		_ = logger.Log(log.Message("message sent"), log.KV{
			K: "message",
			V: msg,
		})
	}
}

func etcdClient(dialTimeout time.Duration, addr string) (client.Client, error) {
	return client.New(client.Config{
		Endpoints: []string{addr},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   dialTimeout,
				KeepAlive: time.Minute,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	})
}
