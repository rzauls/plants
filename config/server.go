package config

import "log/slog"

// env variable keys
const ENV_API_HOST = "API_HOST"
const ENV_API_PORT = "API_PORT"
const ENV_API_ROOT_PREFIX = "API_ROOT_PREFIX"

// default values
const API_DEFAULT_HOST = "localhost"
const API_DEFAULT_PORT = "8080"
const API_DEFAULT_ROOT_PREFIX = "/api/v1"

type Server struct {
	Host       string
	Port       string
	RootPrefix string
}

func FromEnv(logger *slog.Logger, getenv func(string) string) Server {
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

	root := getenv(ENV_API_ROOT_PREFIX)
	if root == "" {
		fallbackWarning(logger, "root prefix", API_DEFAULT_ROOT_PREFIX)
		root = API_DEFAULT_ROOT_PREFIX
	}

	return Server{
		Host:       host,
		Port:       port,
		RootPrefix: root,
	}

}

func fallbackWarning(logger *slog.Logger, fieldName, value string) {
	logger.Warn("using fallback " + fieldName + " value: " + value)
}

func NewDefaultServer() Server {
	return Server{
		Host:       API_DEFAULT_HOST,
		Port:       API_DEFAULT_PORT,
		RootPrefix: API_DEFAULT_ROOT_PREFIX,
	}
}
