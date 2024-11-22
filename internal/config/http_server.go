package config

import "github.com/orungrau/em_song_library/pkg/transport"

type HttpServerConfig struct {
	Address string `env:"SERVER_ADDRESS" env-default:"0.0.0.0:8080"`
}

func NewHttpServerConfig() transport.HTTPServerConfig {
	return &HttpServerConfig{}
}

func (h *HttpServerConfig) GetAddress() string {
	return h.Address
}
