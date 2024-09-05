package configs

import (
	"fmt"
	"os"

	"github.com/spf13/cast"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	ServiceName string
	Environment string
	LoggerLevel string

	AuthServiceGrpcHost string
	AuthServiceGrpcPort string
	Email               string
	Password            string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found", err)
	}

	cfg := Config{}

	cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	cfg.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", 5432))
	cfg.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	cfg.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "0509"))
	cfg.PostgresDB = cast.ToString(getOrReturnDefault("POSTGRES_DB", "tasks"))

	cfg.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "auth_service1"))
	cfg.LoggerLevel = cast.ToString(getOrReturnDefault("LOGGER_LEVEL", "debug"))

	cfg.AuthServiceGrpcHost = cast.ToString(getOrReturnDefault("AUTH_SERVICE_GRPC_HOST", "localhost"))
	cfg.AuthServiceGrpcPort = cast.ToString(getOrReturnDefault("AUTH_SERVICE_GRPC_PORT", ":8081"))

	cfg.Email = cast.ToString(getOrReturnDefault("EMAIL", "boburbekyusupov42@gmail.com"))
	cfg.Password = cast.ToString(getOrReturnDefault("Password", "xmbj bguc xean hsii"))

	return cfg
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}
