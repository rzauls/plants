package httpd

type httpError struct {
	Message string `json:"message"`
}

func newHttpError(err error) httpError {
	return httpError{
		Message: err.Error(),
	}
}

type validationError struct {
	Message  string            `json:"message"`
	Problems map[string]string `json:"errors"`
}

func newValidationError(msg string, problems map[string]string) validationError {
	return validationError{
		Message:  msg,
		Problems: problems,
	}
}
