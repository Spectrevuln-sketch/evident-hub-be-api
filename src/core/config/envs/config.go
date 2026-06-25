package envs

import (
	"sync"
)

type Config struct {
	AppConfig *appConfig
	DbConfig  *databaseConfig
}

var (
	cfg  *Config
	once sync.Once
)

func Get() *Config {
	once.Do(func() {
		cfg = &Config{
			AppConfig: newAppConfig(),
			DbConfig:  newDbConfig(),
		}
	})

	return cfg
}
