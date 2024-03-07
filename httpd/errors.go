package httpd

type httpError struct {
	Message string `json:"message"`
}

func newHttpError(err error) httpError {
	return httpError{
		Message: err.Error(),
	}
}
