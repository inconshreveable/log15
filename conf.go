package log15

import(
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

// config file will be parsed to be Configuration struct
type Configuration struct{
    handlers map[string]*CfgHandler
    loggers map[string]*CfgLogger
}

type CfgHandler struct{
    handler_type string
    path string
}

type CfgLogger struct{
    level string
    handler_name string
}


// read from config file to make Configuration struct.
func readConfig() (*Configuration, error){
    conf := make(map[interface{}]interface{})
    buf, err := ioutil.ReadFile("./log15.yml")
    if nil == err {
        err = yaml.Unmarshal(buf, &conf)
        if nil == err{
            return makeConfiguration(conf)
        }
        return nil, err
    }
    return nil, err
}

func newConfiguration() *Configuration {
    out := new(Configuration)
    out.handlers = make(map[string]*CfgHandler)
    out.loggers = make(map[string]*CfgLogger)
    return out
}

func makeConfiguration(confMap map[interface{}]interface{}) (*Configuration, error) {
    cmloggers := confMap["loggers"].(map[interface{}]interface{})
    cmhandlers := confMap["handlers"].(map[interface{}]interface{})

    out := newConfiguration()
    for lgName, cmlogger := range cmloggers {
        logger := new(CfgLogger) 
        for k, v := range cmlogger.(map[interface{}]interface{}) {
            if k.(string) == "level"{
                logger.level = v.(string)
            }else if k.(string) == "handler_name"{
                logger.handler_name = v.(string)
            }
        }
        out.loggers[lgName.(string)] = logger
    }
    for hdName, cmhandler := range cmhandlers {
        handler := new(CfgHandler) 
        for k, v := range cmhandler.(map[interface{}]interface{}) {
            if k.(string) == "handler_type" {
                handler.handler_type = v.(string)
            }else if k.(string) == "path"{
                handler.path = v.(string)
            }
        }
        out.handlers[hdName.(string)] = handler
    }
    return out, nil
}
