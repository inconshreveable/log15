# About 

Package config enables configuration of log15 from some arbitrary maps. It does not read any config files, you can use 
your favourite way of configuring your application.
This is a deliberate design decition in order to decouple reading of a configuration from the act of instantiating
Handlers and Loggers from this configuration.

# Usage

See this example. Detailled examples for all Handlers can be found in the tests.

```go
func getMapFromConfiguration(config string) (map[string]interface{}, error) {
	configMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(config), &configMap)
	if err != nil {
		return nil, err
	}
	return configMap, err
}

func Example() {
	var exampleConfiguration = `
  # default for all handlers
  level: INFO
  extra:
      mark: test
      user: alice

  handlers:
    - kind: stdout  # determines the Handler used
      format: terminal

    - kind: stderr
      format: json
      level: warn	# don't show

    - kind: stdout
      format: logfmt
      level: debug
`

	configMap, err := getMapFromConfiguration(exampleConfiguration)
	if err != nil {
		panic(err)
	}

	log, err := config.Logger(configMap)
	if err != nil {
		panic(err)
	}

	log.Info("Hello, world!")
	log.Debug("user1", "user", "bob")

	l1 := log.New("user", "carol") // issue in log15! won't override, but use both!
	l1.Debug("about user")
}
```

The kinds for the Handlers are:

* stdout: StreamHandler to os.Stdout
    - format: see Formats below
    
* stderr: StreamHandler to os.Stderr
    - format: see Formats below
    
* file: Fileandler to a file
    - format: see Formats below
    - path: path to file

not yet implemented:

* net
* syslog
* syslog_net

## Format: 

one of "terminal","json","logfmt"

# Add new Handlers

This is accomplished with the help of (mapstructure-hooks)[https://]github.com/gernoteger/mapstructure-hooks]. You need:

1. A struct that will hold your config and implements HandlerConfig. Use LevelHandlerConfig for proper level handling.

```go
type AwesomeHandlerConfig struct {
        LevelHandlerConfig `mapstructure:",squash"`
        AwesomeData string // an example for your data
}
```
Take care that alle fields are exported, otherwise reflection in mapstruct will fail!

2. A function to create an empty config to be filled by mapstruct:

```go
func NewAwesomeHandlerConfig() interface{}{
        return &AwesomeHandlerConfig{}
}
```

3. Register from init() function

add to config/handler.go:

```go
hooks.Register(HandlerConfigType, "gelf", NewGelfConfig)
```
