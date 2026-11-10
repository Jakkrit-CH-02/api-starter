// Package env loads the service configuration from environment variables,
// optionally read from a .env file.
package env

import (
	"os"
	"reflect"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	// mountedEnvFile is used when it exists (e.g. mounted into a container);
	// otherwise .env in the working directory is used.
	mountedEnvFile = "/configs/.env"

	defaultPort = "3005"
)

// Config holds every setting of the service. The `env` tag is the name of the
// environment variable the field is read from.
type Config struct {
	ProductCode string `env:"PRODUCT_CODE"`
	ModuleName  string `env:"MODULE_NAME"`
	Build       string `env:"BUILD"`
	Release     string `env:"RELEASE"`
	Port        string `env:"PORT"`

	DBServer   string `env:"DB_SERVER"`
	DBPort     string `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBFedAuth  string `env:"DB_FED_AUTH"`
}

var (
	instance *Config
	once     sync.Once
)

// Env returns the configuration, loading it on the first call.
func Env() *Config {
	once.Do(func() {
		instance = load()
	})
	return instance
}

func load() *Config {
	if _, err := os.Stat(mountedEnvFile); err == nil {
		if err := godotenv.Load(mountedEnvFile); err != nil {
			log.Fatal().Err(err).Msgf("cannot load %s", mountedEnvFile)
		}
	} else if err := godotenv.Load(); err != nil {
		// Not fatal: the variables may be set directly in the environment.
		log.Debug().Err(err).Msg("no .env file loaded, using process environment")
	}

	cfg := &Config{}
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if key := t.Field(i).Tag.Get("env"); key != "" {
			v.Field(i).SetString(os.Getenv(key))
		}
	}

	if cfg.Port == "" {
		cfg.Port = defaultPort
	}
	return cfg
}
