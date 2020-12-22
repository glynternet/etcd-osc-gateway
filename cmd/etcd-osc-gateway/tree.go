package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/glynternet/etcd-osc-gateway/pkg/etcd"
	"github.com/glynternet/etcd-osc-gateway/pkg/osc"
	osc2 "github.com/glynternet/go-osc/osc"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.etcd.io/etcd/clientv3"
)

const (
	defaultDialTimeout = time.Second
	defaultReadTimeout = 250 * time.Millisecond
	requestTimeout     = 2 * time.Second
)

func buildCmdTree(logger log.Logger, _ io.Writer, rootCmd *cobra.Command) {
	var listenHost string
	var listenPort uint
	var etcdCfg etcdDialConfig

	rootCmd.RunE = func(_ *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		return run(ctx, logger, listenHost, listenPort, etcdCfg)
	}

	rootCmd.Flags().StringVar(&listenHost, "listen-host", "127.0.0.1", "host address to listen on")
	rootCmd.Flags().UintVar(&listenPort, "listen-port", 9000, "host post to listen on")
	rootCmd.Flags().StringVar(&etcdCfg.scheme, "etcd-scheme", "http", "etcd scheme")
	rootCmd.Flags().StringVar(&etcdCfg.host, "etcd-host", "127.0.0.1", "etcd host")
	rootCmd.Flags().UintVar(&etcdCfg.port, "etcd-port", 2379, "etcd port")
}

func run(_ context.Context, logger log.Logger, listenHost string, listenPort uint, etcdCfg etcdDialConfig) error {
	etcdDialAddr := etcdCfg.dialAddress()
	cli, err := etcdClient(defaultDialTimeout, etcdDialAddr)
	if err != nil {
		return errors.Wrap(err, "creating client")
	}
	if err := logger.Log(log.Message("Client created at address:%s"),
		log.KV{
			K: "dialAddress",
			V: etcdDialAddr,
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
			KeyValuePutter: etcd.Client{KV: clientv3.NewKV(cli)},
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
			K: "oscMessage",
			V: msg,
		})
	}
}

func etcdClient(dialTimeout time.Duration, addr string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:            []string{addr},
		DialTimeout:          dialTimeout,
		DialKeepAliveTime:    5 * time.Second,
		DialKeepAliveTimeout: time.Minute,
	})
}
