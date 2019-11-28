package cipherassets_core

import (
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

type Config struct {
	Database *databaseConfig
	Debug    bool
}

type databaseConfig struct {
	driver string
	host   string
	name   string
	pass   string
	port   int
	user   string
}

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

	c := &Config{
		Debug:    getEnvBool("DEBUG", false),
		Database: &dbConf,
	}

	return c, nil
}
