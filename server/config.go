package main

import (
	"errors"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
    AmqpUrl string `toml:"url"`
    Port string `toml:"port"`
    FilterPattern []string `toml:"filter"`
}

func ParseConfigFile(file []byte) (ServerConfig, error) {
    sc := ServerConfig{}
    _, err := toml.Decode(string(file), &sc)
    if err != nil {
        return ServerConfig{}, err
    }

    if sc.AmqpUrl == "" {
        return ServerConfig{}, errors.New("Empty url")
    }

    if sc.Port == "" {
        return ServerConfig{}, errors.New("Empty port")
    }

    if len(sc.FilterPattern) == 0 {
        return ServerConfig{}, errors.New("Empty filters")
    }

    return sc, nil
}
