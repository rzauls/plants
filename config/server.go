package config

type Server struct {
	Host       string
	Port       string
	RootPrefix string
}

func NewDefaultServer() Server {
	return Server{
		Host:       "localhost",
		Port:       "8080",
		RootPrefix: "/api/v1",
	}
}
