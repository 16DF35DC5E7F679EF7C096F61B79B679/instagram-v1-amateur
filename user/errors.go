package user

type CreationValidationError struct {
	InvalidField string
	ErrorReason string
}

func (creationValidationError *CreationValidationError) Error() string {
	return "Error in User Creation: " +
		creationValidationError.InvalidField + " is invalid: " + creationValidationError.ErrorReason
}