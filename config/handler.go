package config

import (
	"fmt"
	"os"
	"reflect"

	"net/url"

	"github.com/gernoteger/mapstructure-hooks"
	"github.com/inconshreveable/log15"
)

// Just the selector for the Handler!
type Handler int

// GLEF is not a  format, it's a handler!!
const (
	HandlerStdout Handler = iota
	HandlerStderr
)

func init() {
	Register()
}

// HandlerConfig will create Handlers from a config
type HandlerConfig interface {
	NewHandler() (log15.Handler, error)
	GetLevel() string
}

// use for registry functions
var HandlerConfigType = reflect.TypeOf((*HandlerConfig)(nil)).Elem()

// registers all handlers
func Register() {
	hooks.RegisterInterface(HandlerConfigType, "kind")

	hooks.Register(HandlerConfigType, "stdout", NewStdoutConfig)
	hooks.Register(HandlerConfigType, "stderr", NewStderrConfig)
	hooks.Register(HandlerConfigType, "file", NewFileConfig)
	hooks.Register(HandlerConfigType, "gelf", NewGelfConfig)
	hooks.Register(HandlerConfigType, "net", NewNetConfig)

	hooks.Register(HandlerConfigType, "buffer", NewBufferConfig)
	hooks.Register(HandlerConfigType, "multi", NewMultiConfig)
	hooks.Register(HandlerConfigType, "filter", NewMatchFilterConfig)
	hooks.Register(HandlerConfigType, "failover", NewFailoverConfig)

}

type LevelHandlerConfig struct {
	Level string
}

func (c *LevelHandlerConfig) GetLevel() string {
	return c.Level
}

type StreamConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Handler            Handler // for differentiation of stdion, stdout
	Format             Fmt
}

func (c *StreamConfig) NewHandler() (log15.Handler, error) {
	var f *os.File
	switch c.Handler {
	case HandlerStdout:
		f = os.Stdout
	case HandlerStderr:
		f = os.Stderr
	default:
		return nil, fmt.Errorf("invalid handler: %v", c.Handler)
	}

	log := log15.StreamHandler(f, c.Format.NewFormat())
	//TODO: use level

	return log, nil
}
func (c *StreamConfig) GetLevel() string {
	return c.Level
}

// return a ConsoleConfig with default values
func NewStreamConfig() *StreamConfig {
	return &StreamConfig{} //File: "stderr"
}

func NewStdoutConfig() interface{} {
	return &StreamConfig{Handler: HandlerStdout}
}
func NewStderrConfig() interface{} {
	return &StreamConfig{Handler: HandlerStderr}
}

type FileConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Path               string
	Format             Fmt
}

func NewFileConfig() interface{} {
	return &FileConfig{}
}

func (c *FileConfig) NewHandler() (log15.Handler, error) {
	h, err := log15.FileHandler(c.Path, c.Format.NewFormat())
	if err != nil {
		return nil, err
	}
	return h, nil
}

type GelfConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Address            string
}

// make sure its's the right interface
var _ HandlerConfig = (*GelfConfig)(nil)

func NewGelfConfig() interface{} {
	return &GelfConfig{}
}

func (c *GelfConfig) NewHandler() (log15.Handler, error) {
	h, err := log15.GelfHandler(c.Address)
	return h, err
}

type NetConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Format             Fmt
	URL                string
}

// make sure its's the right interface
var _ HandlerConfig = (*NetConfig)(nil)

func NewNetConfig() interface{} {
	return &NetConfig{}
}

func (c *NetConfig) NewHandler() (log15.Handler, error) {
	u, err := url.Parse(c.URL)
	if err != nil {
		return nil, err
	}

	h, err := log15.NetHandler(u.Scheme, u.Host, c.Format.NewFormat())
	if err != nil {
		return nil, err
	}
	return h, err
}

// BufferConfig is a buffered handler
type BufferConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Handler            HandlerConfig
	BufSize            int
}

// make sure its's the right interface
var _ HandlerConfig = (*BufferConfig)(nil)

func NewBufferConfig() interface{} {
	return &BufferConfig{BufSize: 10}
}

func (c *BufferConfig) NewHandler() (log15.Handler, error) {

	h, err := c.Handler.NewHandler()
	if err != nil {
		return nil, err
	}
	return log15.BufferedHandler(c.BufSize, h), nil
}

// ---- MultiHandler

// MultiHandler fans out ot all handlers
type MultiConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Handlers           []HandlerConfig
}

// make sure its's the right interface
var _ HandlerConfig = (*MultiConfig)(nil)

func NewMultiConfig() interface{} {
	return &MultiConfig{}
}

// handlers creates handles form []HandlersConfig
func handlers(c []HandlerConfig) ([]log15.Handler, error) {
	hh := make([]log15.Handler, len(c))
	for i, hc := range c {
		var err error
		hh[i], err = hc.NewHandler()
		if err != nil {
			return nil, err
		}
	}
	return hh, nil
}

func (c *MultiConfig) NewHandler() (log15.Handler, error) {
	// make 'em all
	hh, err := handlers(c.Handlers)
	if err != nil {
		return nil, err
	}
	return log15.MultiHandler(hh...), nil
}

//-------------- MatchFilterHandler

// MatchFilterHandler onyl fires if field matches value
type MatchFilterConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Handler            HandlerConfig
	Key                string
	Value              interface{}
}

// make sure its's the right interface
var _ HandlerConfig = (*MatchFilterConfig)(nil)

func NewMatchFilterConfig() interface{} {
	return &MatchFilterConfig{}
}

func (c *MatchFilterConfig) NewHandler() (log15.Handler, error) {
	h, err := c.Handler.NewHandler()
	if err != nil {
		return nil, err
	}
	return log15.MatchFilterHandler(c.Key, c.Value, h), nil
}

// ---- FailoverHandler

// FailoverConfig configure FailoverHandler
type FailoverConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Handlers           []HandlerConfig
}

// make sure its's the right interface
var _ HandlerConfig = (*FailoverConfig)(nil)

func NewFailoverConfig() interface{} {
	return &FailoverConfig{}
}

func (c *FailoverConfig) NewHandler() (log15.Handler, error) {
	// make 'em all
	hh, err := handlers(c.Handlers)
	if err != nil {
		return nil, err
	}
	return log15.FailoverHandler(hh...), nil
}
