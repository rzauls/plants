package config

type Server struct {
	Host string
	Port int
}

func NewDefaultServer() Server {
	return Server{
		Host: "localhost",
		Port: 8080,
	}
}
