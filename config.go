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
	App      *appConfig
	Database *databaseConfig
	JWT      *jwtConfig
}

type appConfig struct {
	Name  string
	Debug bool
}

type databaseConfig struct {
	Driver string
	Dsn    string
	Host   string
	Name   string
	Pass   string
	Port   int
	User   string
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
	valid := len(c.Driver) != 0 &&
		len(c.Host) != 0 &&
		len(c.Name) != 0 &&
		len(c.Pass) != 0 &&
		c.Port != 0 &&
		len(c.User) != 0
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
		Driver: getEnvString("DB_DRIVER", ""),
		Host:   getEnvString("DB_HOST", ""),
		Name:   getEnvString("DB_NAME", ""),
		Pass:   getEnvString("DB_PASSWORD", ""),
		Port:   getEnvInt("DB_PORT", 0),
		User:   getEnvString("DB_USER", ""),
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

	appConf := appConfig{
		Name:  getEnvString("NAME", "Cipher Assets"),
		Debug: getEnvBool("DEBUG", false),
	}

	c := &Config{
		App:      &appConf,
		Database: &dbConf,
		JWT:      &jwtConf,
	}

	return c, nil
}
