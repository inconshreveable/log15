package log15

import (
    "fmt"
)

type Loggers map[string]*logger

var (
    // store logger which is initialized from configuration file.
    loggers Loggers 
)

func initByConf() {
    loggers = Loggers{"root":root}
    makeLoggersFromConf()
}

// get specific logger by name.
// when you get logger by this method,you can use logger with the former way.
// return root if the name is not defined in the configuration file.
func GetLogger(name string) (*logger, error){
    if nil == loggers {
        initByConf()
    }
    if v, ok := loggers[name]; ok{
        return v, nil
    }else{
        return loggers["root"], nil
    }
}

func makeHandler(conf *Configuration, name string) (Handler, error){
    cfHandler, ok := conf.handlers[name]
    if !ok {
        return nil, fmt.Errorf("no handler named %s defined in configuraiton file", name)
    }
    var (
        err error
        handler Handler
    )
    switch cfHandler.handler_type{
    case "file":
        handler, err = FileHandler(cfHandler.path, LogfmtFormat())
        if err != nil {
            return nil, err
        }
    case "stdout":
        handler = StdoutHandler
    }
    return handler, err
}

func makeLogger(conf *Configuration, name string) (*logger, error){
    handler, err := makeHandler(conf, conf.loggers[name].handler_name)
    if err != nil {
        return nil, err
    }
    lg := &logger{[]interface{}{}, new(swapHandler)} 
    lg.SetHandler(handler)
    return lg, nil
}

func makeLoggersFromConf() {
    conf, err := readConfig()
    if err != nil {
        fmt.Println(err)
        return 
    }
    for lgname, _ := range conf.loggers {
        loggers[lgname], err = makeLogger(conf, lgname)
        if err != nil {
            continue
        }
    }
}
