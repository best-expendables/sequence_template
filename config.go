package sequencetemplate

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

// PostgresConfig config for gorm postgres connection
type PostgresConfig struct {
	Host   string `envconfig:"POSTGRES_HOST" required:"true"`
	Port   int    `envconfig:"POSTGRES_PORT" required:"true"`
	DbName string `envconfig:"POSTGRES_DB" required:"true"`
	User   string `envconfig:"POSTGRES_USER" required:"true"`
	Pass   string `envconfig:"POSTGRES_PASS" required:"true"`
}

func (conf *PostgresConfig) getConnectionURL() string {
	return fmt.Sprintf(
		"host=%s user=%s port=%d dbname=%s sslmode=disable password=%s",
		conf.Host,
		conf.User,
		conf.Port,
		conf.DbName,
		conf.Pass,
	)
}

func getAppConfigFromEnv() *PostgresConfig {
	var conf PostgresConfig
	envconfig.MustProcess("", &conf)
	return &conf
}
