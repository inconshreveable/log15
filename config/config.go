package config

import (
	"fmt"
	"github.com/inconshreveable/log15"
	"strings"
	"github.com/gernoteger/mapstructure-hooks"
)

// LoggerConfig is the central configuration that will be populated from logfiles by various Method
// One LoggerConfig will produce one logger
type LoggerConfig struct {
	Level    string
	Handlers []HandlerConfig

	// extra fields to be added
	//Extra map[string]interface{}
	Extra log15.Ctx
}

// NewLogger produces a new logger from a configuration
func (c *LoggerConfig) NewLogger() (log15.Logger, error) {
	var handlers []log15.Handler

	if c.Level == "" {
		c.Level="info"
	}
	for _, hc := range c.Handlers {
		if hc == nil {
			return nil, fmt.Errorf("nil handler")
		}
		h, err := hc.NewHandler()
		if err != nil {
			return nil, err
		}

		//set log level
		//TODO: use level: how to get root level??
		l:=c.Level
		if hc.GetLevel()!="" {
			l=hc.GetLevel()
		}

		lvl, err := log15.LvlFromString(strings.ToLower(l))
		if err != nil {
			//TODO: better explanation!
			return nil,err
		}
		h = log15.LvlFilterHandler(lvl, h)

		handlers = append(handlers, h)
	}

	hall := log15.MultiHandler(handlers...)

	l := log15.New(c.Extra)
	l.SetHandler(hall)
	return l, nil
}

// LoggerConfig creates a new config from map data
func NewLoggerConfig(configMap map[string]interface{}) (*LoggerConfig,error) {
	c := LoggerConfig{}
	err := hooks.Decode(configMap, &c)
	if err != nil {
		return nil, err
	}
	return &c,nil
}

// Logger creates a new Logger from a configuration map
func Logger(config map[string]interface{}) (log15.Logger,error) {

	configMap,err:=NewLoggerConfig(config)
	if err != nil {
		return nil,err
	}
	// das gehört zusammen!!
	c := LoggerConfig{}
	err = hooks.Decode(configMap, &c)
	if err != nil {
		return nil, err
	}
	return c.NewLogger()
}