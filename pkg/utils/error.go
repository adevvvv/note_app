package utils

// HTTPError представляет ошибку HTTP с сообщением и кодом состояния.
type HTTPError struct {
	Message string
	Code    int
}

// Error возвращает текстовое представление ошибки.
func (e *HTTPError) Error() string {
	return e.Message
}
