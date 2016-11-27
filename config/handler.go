package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/inconshreveable/log15"
	"github.com/gernoteger/mapstructure-hooks"

)

// Just the selector for the Handler!
type Handler int

// GLEF is not a  format, it's a handler!!
const (
	HandlerStdout Handler = iota
	HandlerStderr
	HandlerFile
	HandlerNet
	HandlerSyslog
	HandlerSyslogNet
	HandlerGraylog
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
	hooks.RegisterInterface(HandlerConfigType,"kind")

	hooks.Register(HandlerConfigType, "stdout", NewStdoutConfig)
	hooks.Register(HandlerConfigType, "stderr", NewStderrConfig)
	hooks.Register(HandlerConfigType, "file", NewFileConfig)
	hooks.Register(HandlerConfigType, "gelf", NewGelfConfig)
}


type LevelHandlerConfig struct {
	Level  string
}

func (c * LevelHandlerConfig) GetLevel() string {
	return c.Level
}

type StreamConfig struct {
	Handler Handler // make it easy..
	Format  Fmt
	Level   string
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

	log:=log15.StreamHandler(f,c.Format.NewFormat())
	//TODO: use level

	return log, nil
}
func (c * StreamConfig) GetLevel() string {
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
	Path	string
	Format  Fmt
}

func NewFileConfig() interface {} {
	return &FileConfig{}
}


func (c *FileConfig) NewHandler() (log15.Handler, error) {
	h,err:=log15.FileHandler(c.Path,c.Format.NewFormat())
	if err!=nil{
		return nil,err
	}
	return h,nil
}


type GelfConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Address string
}

// make sure its's the right interface
var _ HandlerConfig = (*GelfConfig)(nil)


func NewGelfConfig() interface{} {
	return &GelfConfig{}
}

func (c * GelfConfig) NewHandler() (log15.Handler, error) {
	h,err:=log15.GelfHandler(c.Address)
	return h,err
}

