package httpd

type httpError struct {
    Code string `json:"code"`
    Message string `json:"message"`
}
