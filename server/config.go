package main

import (
	"fmt"
	"regexp"

	"github.com/BurntSushi/toml"
)

type serverConfig struct {
	Url           string   `toml:"url"`
	Port          string   `toml:"port"`
	FilterPattern []string `toml:"filter"`
}

func parseConfig(file []byte) (serverConfig, error) {
	sc := serverConfig{}
	_, err := toml.Decode(string(file), &sc)
	if err != nil {
		return serverConfig{}, err
	}

	return sc, nil
}

func (sc serverConfig) genFilters() ([]*regexp.Regexp, error) {
	var filters []*regexp.Regexp
	for _, v := range sc.FilterPattern {
		regex, err := regexp.Compile(v)
		if err != nil {
			return []*regexp.Regexp{}, err
		}
		filters = append(filters, regex)
	}

	return filters, nil
}

func (sc serverConfig) String() string {
	return fmt.Sprintf(
		`URL      = %v
PORT     = %v
FILTER   = %v`,
		sc.Url,
		sc.Port,
		sc.FilterPattern,
	)
}
