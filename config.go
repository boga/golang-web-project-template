package cipherassets_core

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

type Config struct {
	Database *databaseConfig
	Debug    bool
	JWT      *jwtConfig
}

type databaseConfig struct {
	driver string
	dsn    string
	host   string
	name   string
	pass   string
	port   int
	user   string
}

type jwtConfig struct {
	AuthTokenTTL *time.Duration
	Secret       string
}

func (c jwtConfig) Valid() error {
	if c.Secret == "" || c.AuthTokenTTL == nil {
		return fmt.Errorf("jwt config is not valid, some options missing")
	}

	return nil
}

// func (c databaseConfig) GetDSN() error {
// 	var dsn string
// 	if c.dsn == "" {
// 		dsn = fmt.Sprintf("%s://%s:%s@%s:%d/%s", c.driver,
// 			c.user, c.pass, c.host, c.port, c.name)
// 	}
// 	if c.Valid() == nil {
// 		// "test:test@(localhost:3306)/test"
// 		dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s", c.driver,
// 			c.user, c.pass, c.host, c.port, c.name)
// 	}
//
// 	return dsn
// }

func (c databaseConfig) Valid() error {
	valid := len(c.driver) != 0 &&
		len(c.host) != 0 &&
		len(c.name) != 0 &&
		len(c.pass) != 0 &&
		c.port != 0 &&
		len(c.user) != 0
	if !valid {
		return fmt.Errorf("database config is not valid, some options missing")
	}
	return nil
}

func getEnvString(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value, err := strconv.Atoi(getEnvString(key, "non-integer-string"))
	if err != nil {
		value = fallback
	}

	return value
}

func getEnvBool(key string, fallback bool) bool {
	valueStr := getEnvString(key, strconv.FormatBool(fallback))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		value = fallback
		log.Fatalf("Option %s has erroneus value \"%s\"", key, valueStr)
	}

	return value
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("got error while loading .env file: %s", err.Error())
	}

	dbConf := databaseConfig{
		driver: getEnvString("DB_DRIVER", ""),
		host:   getEnvString("DB_HOST", ""),
		name:   getEnvString("DB_NAME", ""),
		pass:   getEnvString("DB_PASSWORD", ""),
		port:   getEnvInt("DB_PORT", 0),
		user:   getEnvString("DB_USER", ""),
	}
	if err := dbConf.Valid(); err != nil {
		return nil, err
	}

	AuthTokenTTLStr := getEnvString("JWT_AUTH_TTL", "") // 30 days
	AuthTokenTTL, err := time.ParseDuration(AuthTokenTTLStr)
	if err != nil {
		panic("Cant't parse duration value " + AuthTokenTTLStr)
	}
	jwtConf := jwtConfig{
		AuthTokenTTL: &AuthTokenTTL,
		Secret:       getEnvString("JWT_SECRET", ""),
	}
	if err := jwtConf.Valid(); err != nil {
		return nil, err
	}

	c := &Config{
		Debug:    getEnvBool("DEBUG", false),
		Database: &dbConf,
		JWT:      &jwtConf,
	}

	return c, nil
}
