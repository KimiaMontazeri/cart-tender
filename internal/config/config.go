package config

import (
	"log"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

const (
	// Prefix indicates environment variables prefix.
	Prefix = "carttender_"
)

type (
	// Config holds all configurations.
	Config struct {
		API      API          `koanf:"api"`
		Postgres PostgresDB   `koanf:"postgres"`
		JWT      JsonWebToken `koanf:"jwt"`
	}

	API struct {
		Port int `koanf:"port"`
	}

	PostgresDB struct {
		Host           string        `koanf:"host"`
		Port           int           `koanf:"port"`
		Username       string        `koanf:"username"`
		Password       string        `koanf:"password"`
		DBName         string        `koanf:"dbname"`
		ConnectTimeout time.Duration `koanf:"connect-timeout"`
	}

	JsonWebToken struct {
		Expiration        time.Duration `koanf:"expiration"`
		ServiceExpiration time.Duration `koanf:"service-expiration"`
		Secret            string        `koanf:"secret"`
	}
)

// New reads configuration with viper.
func New() Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		log.Printf("error loading config.yml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "_", ".")
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	log.Printf("following configuration is loaded:\n%+v", instance)

	return instance
}
