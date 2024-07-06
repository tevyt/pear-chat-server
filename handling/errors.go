package handling

type AuthenticationError struct {
	message string
}

func NewAuthenticationError(message string) AuthenticationError {
	return AuthenticationError{message: message}
}
func (err AuthenticationError) Error() string {
	return err.message
}
