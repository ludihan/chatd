package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
    Url string `toml:"url"`
    Port string `toml:"port"`
    FilterPattern []string `toml:"filter"`
    Exchange string `toml:"exchange"`
}

func ParseConfig(file []byte) (ServerConfig, error) {
    sc := ServerConfig{}
    _, err := toml.Decode(string(file), &sc)
    if err != nil {
        return ServerConfig{}, err
    }

    return sc, nil
}

func (sc ServerConfig) String() string {
    return fmt.Sprintf(
`URL      = %v
PORT     = %v
FILTER   = %v
EXCHANGE = %v`,
    sc.Url,
    sc.Port,
    sc.FilterPattern,
    sc.Exchange,
)
}
