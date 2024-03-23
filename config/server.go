package config

import "log/slog"

// env variable keys
const ENV_API_HOST = "API_HOST"
const ENV_API_PORT = "API_PORT"

// default values
const API_DEFAULT_HOST = "localhost"
const API_DEFAULT_PORT = "8080"

type Server struct {
	Host string
	Port string
}

func FromEnv(getenv func(string) string) Server {
	logger := slog.Default()
	host := getenv(ENV_API_HOST)
	if host == "" {
		fallbackWarning(logger, "host", API_DEFAULT_HOST)
		host = API_DEFAULT_HOST
	}

	port := getenv(ENV_API_PORT)
	if port == "" {

		fallbackWarning(logger, "port", API_DEFAULT_PORT)
		port = API_DEFAULT_PORT
	}

	return Server{
		Host: host,
		Port: port,
	}

}

func fallbackWarning(logger *slog.Logger, fieldName, value string) {
	logger.Warn("using fallback " + fieldName + " value: " + value)
}

func NewDefaultServer() Server {
	return Server{
		Host: API_DEFAULT_HOST,
		Port: API_DEFAULT_PORT,
	}
}
