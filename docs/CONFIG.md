## Config

### What is it?

Viper is a complete configuration solution for Go applications including [12-Factor](https://12factor.net/#the_twelve_factors) apps. It is designed to work within an application, and can handle all types of configuration needs and formats. It supports:

- Load environment variables
- Find load, umarshal config file (JSON, TOML, YAML, ENV, INI)
- Read config from environment variables or flags (Override existing values, set default values)
- Read config from remote system (etcd, Consul)
- Live watching and writing config file (Re-read changed files, save any modifications)

Read more about this from [here](https://github.com/spf13/viper?tab=readme-ov-file).

Install it using
```sh
go get github.com/spf13/viper
```

### Loading env variables

Create a file under package eg: `util` that loads the env variables.

```go
package util

import "github.com/spf13/viper"

// Store all configuration of the application.
// Values are read by viper from a config file or `env` variables.
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// Read configurations from file or `env` variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json, xml, toml

	// Load and overrid config values if found
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}
```

All the env variables are not loaded inside the `config` struct. Now use the as follows:

```go
import (
	"simple-bank/util"
)

config, err := util.LoadConfig(".")
conn, err := sql.Open(config.DBDriver, config.DBSource)
```

### Override values

```sh
SERVER_ADDRESS=0.0.0.0:8081 make server
```