package config_test

import (
	"github.com/inconshreveable/log15/config"
	"gopkg.in/yaml.v2"
)

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
    - kind: stdout
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

	// disabling output below for tests by immediately prepending this line since dates will never be right. to execute, just insert blank line.
	// Output:
	// INFO[11-30|11:37:20] Hello, world!                            mark=test user=alice
	// t=2016-11-30T11:37:20+0100 lvl=info msg="Hello, world!" mark=test user=alice
	// t=2016-11-30T11:37:20+0100 lvl=dbug msg=user1 mark=test user=alice user=bob
	// t=2016-11-30T11:37:20+0100 lvl=dbug msg="about user" mark=test user=alice user=carol

}
