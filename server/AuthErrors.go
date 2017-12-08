package server

const (
	AuthMessageUserNotJoined     = "user has to join first"
	AuthMessageUserAlreadyJoined = "user has already joined"
)

type AuthErrors struct {
	message string
}

func NewAuthError(message string) *AuthErrors {
	return &AuthErrors{message: message}
}

func (e *AuthErrors) Error() string {
	return e.message
}
