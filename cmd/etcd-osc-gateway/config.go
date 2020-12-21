package main

import "fmt"

type etcdDialConfig struct {
	scheme string
	host   string
	port   uint
}

func (cfg etcdDialConfig) dialAddress() string {
	return fmt.Sprintf("%s://%s:%d", cfg.scheme, cfg.host, cfg.port)
}
