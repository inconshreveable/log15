package log15

import(
    "testing"
    "os"
    "io/ioutil"
    "strings"
)

func TestDefaultLogger(t *testing.T){
    lg, err := GetLogger("getdefault")
    if err != nil || nil == lg || lg != root{
        t.Fatalf("Failed to get root logger.")
    }
}

func TestGetFileLogger(t *testing.T){
    err := os.Remove("./test.log")
    if err != nil {
        if !os.IsNotExist(err) {
            t.Fatalf(err.Error())
        }
    }

    // The log file is creaded when the handler is made,
    // so we initialize from config file manually.
    initByConf() 

    lg, err := GetLogger("log_test_log_file")
    if err != nil || nil == lg  || lg == root{
        t.Fatalf("Failed to get file logger.")
    }

    logStr := "write log in file!"
    lg.Info(logStr)
    content,err := ioutil.ReadFile("./test.log")

    if strings.Index(string(content), logStr) < 0 {
        t.Fatalf("Failed to write log in test.log, content:%s", string(content))
    }
}
