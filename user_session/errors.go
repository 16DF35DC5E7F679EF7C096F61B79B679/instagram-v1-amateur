package user_session

type InvalidSessionCreationError struct {
	InvalidField string `json:"invalid_field"`
	ErrorReason string `json:"error_reason"`
}

func (invalidSessionCreationError *InvalidSessionCreationError) Error() string {
	return "Error in session creation: " + " Invalid Field: " + invalidSessionCreationError.InvalidField + " : " + invalidSessionCreationError.ErrorReason
}

type PasswordMismatchError struct {
	RootCause string
}

func (passwordMismatchError *PasswordMismatchError) Error() string {
	return "Invalid Password"
}