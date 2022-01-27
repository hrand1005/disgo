package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type config struct {
	Token string `yaml:"token"`
	//guilds []*string `yaml:"guilds"`
	//features []*string `yaml:"features"`
}

// loadBotConfig parses bot info from the given file
func loadBotConfig(file string) (config, error) {
	c := config{}

	slurp, err := os.ReadFile(file)
	if err != nil {
		return c, err
	}
	if err = yaml.Unmarshal(slurp, &c); err != nil {
		return c, fmt.Errorf("Couldn't unmarshal config: %v", err)
	}

	return c, nil
}
