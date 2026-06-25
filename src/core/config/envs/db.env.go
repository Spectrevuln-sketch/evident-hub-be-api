package envs

import (
	"github.com/kelseyhightower/envconfig"
)

type databaseConfig struct {
	DbName     string `envconfig:"DB_NAME"`
	DbHost     string `envconfig:"DB_HOST"`
	DbUser     string `envconfig:"DB_USER"`
	DbPassword string `envconfig:"DB_PASSWORD"`
	DbPort     int    `envconfig:"DB_PORT"`
	DbMigrate  bool   `envconfig:"DB_MIGRATE"`
	DbSeeder   bool   `envconfig:"DB_SEEDER"`
}

func newDbConfig() *databaseConfig {
	var dbCfg databaseConfig
	envconfig.MustProcess("", &dbCfg)
	return &dbCfg
}
