package helper

type ValidationError struct {
	Errors map[string][]string
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors: make(map[string][]string),
	}
}

func (e *ValidationError) Add(field, message string) {
	e.Errors[field] = append(e.Errors[field], message)
}

func (e *ValidationError) IsEmpty() bool {
	return len(e.Errors) == 0
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

type NotFoundError struct {
	Message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type AuthenticationError struct {
	Message string
	Code    string
}

func (e *AuthenticationError) Error() string {
	return e.Message
}
