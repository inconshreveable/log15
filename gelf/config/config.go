package gelfconfig

import (
	"github.com/gernoteger/mapstructure-hooks"
	"github.com/inconshreveable/log15"

	"github.com/inconshreveable/log15/config"
)

type GelfConfig struct {
	config.LevelHandlerConfig `mapstructure:",squash"`
	Address string
}

// make sure its's the right interface
var _ config.HandlerConfig = (*GelfConfig)(nil)


func NewConfig() interface{} {
	return &GelfConfig{}
}

func (c * GelfConfig) NewHandler() (log15.Handler, error) {
	h,err:=log15.GelfHandler(c.Address)
	return h,err
}

func init() {
	hooks.Register(config.HandlerConfigType, "gelf", NewConfig)
	hooks.Register(config.HandlerConfigType, "graylog", NewConfig)
}
