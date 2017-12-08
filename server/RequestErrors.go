package server

const (
	RequestMessageNoUser = "no user found in the request"
)

type RequestError struct {
	message string
}

func NewRequestError(message string) *RequestError {
	return &RequestError{message:message}
}

func (e* RequestError) Error() string {
	return e.message
}