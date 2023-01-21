# etcd-osc-store

```
etcd-osc-gateway stores OSC messages as key-value pairs in an etcd store.

Any message received is expected to have a single string argument,
which is the value that will be stored in etcd using the message's
address as the key.

If a OSC bundle is received, each message is expected to be in the
same format as described above and each is pushed to etcd.

Usage:
  etcd-osc-gateway [flags]
  etcd-osc-gateway [command]

Available Commands:
  help        Help about any command
  version     show the version of this application

Flags:
      --etcd-host string     etcd host (default "127.0.0.1")
      --etcd-port uint       etcd port (default 2379)
      --etcd-scheme string   etcd scheme (default "http")
  -h, --help                 help for etcd-osc-gateway
      --listen-host string   host address to listen on (default "127.0.0.1")
      --listen-port uint     host post to listen on (default 9000)

Use "etcd-osc-gateway [command] --help" for more information about a command.
```