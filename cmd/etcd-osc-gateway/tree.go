package main

import (
	"fmt"
	"io"
	"time"

	"github.com/glynternet/etcd-osc-gateway/pkg/etcd"
	"github.com/glynternet/etcd-osc-gateway/pkg/osc"
	"github.com/glynternet/pkg/log"
	osc2 "github.com/hypebeast/go-osc/osc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.etcd.io/etcd/clientv3"
)

const (
	defaultDialTimeout = time.Second
	defaultReadTimeout = 250 * time.Millisecond
)

func buildCmdTree(logger log.Logger, _ io.Writer, rootCmd *cobra.Command) {
	var listenHost string
	var listenPort uint
	var etcdCfg etcdDialConfig

	rootCmd.RunE = func(_ *cobra.Command, args []string) error {
		return run(logger, listenHost, listenPort, etcdCfg)
	}

	rootCmd.Flags().StringVar(&listenHost, "listen-host", "127.0.0.1", "host address to listen on")
	rootCmd.Flags().UintVar(&listenPort, "listen-port", 9000, "host post to listen on")
	rootCmd.Flags().StringVar(&etcdCfg.scheme, "etcd-scheme", "http", "etcd scheme")
	rootCmd.Flags().StringVar(&etcdCfg.host, "etcd-host", "127.0.0.1", "etcd host")
	rootCmd.Flags().UintVar(&etcdCfg.port, "etcd-port", 2379, "etcd port")
}

func run(logger log.Logger, listenHost string, listenPort uint, etcdCfg etcdDialConfig) error {
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
		Dispatcher: osc.KeyValueDispatcher{
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
		_ = logger.Log(log.Message("error sending key value pair"), log.ErrorMessage(err))
	}
}

func loggingSuccessHandler(logger log.Logger) func(osc2.Message, string, string) {
	return func(msg osc2.Message, k, v string) {
		_ = logger.Log(log.Message("key value pair sent"),
			log.KV{K: "oscMessage", V: msg},
			log.KV{K: "key", V: k},
			log.KV{K: "value", V: v},
		)
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
