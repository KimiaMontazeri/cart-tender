package config

import "time"

func Default() Config {
	return Config{
		API: API{Port: 65432},
		Postgres: PostgresDB{
			Host:           "localhost",
			Port:           5432,
			DBName:         "cart",
			Username:       "postgres",
			Password:       "postgres",
			ConnectTimeout: 30 * time.Second,
		},
		JWT: JsonWebToken{
			Expiration:        3 * time.Hour,
			ServiceExpiration: 6 * time.Hour,
			Secret:            "secret",
		},
	}
}
