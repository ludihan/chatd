package main

import (
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

    return sc, nil
}
