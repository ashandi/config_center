package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	HttpServerReadTimeout     time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" required:"true"`
	HttpServerWriteTimeout    time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" required:"true"`
	HttpServerShutdownTimeout time.Duration `envconfig:"HTTP_SERVER_SHUTDOWN_TIMEOUT" required:"true"`
	HttpServerAddress         string        `envconfig:"HTTP_SERVER_ADDRESS" required:"true"`

	MysqlDbUser     string `envconfig:"MYSQL_DB_USER" required:"true"`
	MysqlDbPassword string `envconfig:"MYSQL_DB_PASSWORD" required:"true"`
	MysqlDbHost     string `envconfig:"MYSQL_DB_HOST" required:"true"`
	MysqlDbName     string `envconfig:"MYSQL_DB_NAME" required:"true"`

	RedisHost     string        `envconfig:"REDIS_HOST" required:"true"`
	RedisPassword string        `envconfig:"REDIS_PASSWORD" required:"true"`
	RedisDb       int           `envconfig:"REDIS_DB" required:"true"`
	RedisTtl      time.Duration `envconfig:"REDIS_TTL" required:"true"`

	AppVersionRequired         string `envconfig:"APP_VERSION_REQUIRED" required:"true"`
	AppVersionStore            string `envconfig:"APP_VERSION_STORE" required:"true"`
	AppBackendEntrypoint       string `envconfig:"APP_BACKEND_ENTRYPOINT" required:"true"`
	AppNotificationsEntrypoint string `envconfig:"APP_NOTIFICATIONS_ENTRYPOINT" required:"true"`
}

func FromEnv() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}
