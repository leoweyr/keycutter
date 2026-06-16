package fault

// ValidationError represents a structural or integrity failure detected while a
// token is checked against the Keycutter specification.
type ValidationError struct {
	reason string
}

// NewValidationError creates a ValidationError carrying the given human-readable reason.
func NewValidationError(reason string) *ValidationError {
	return &ValidationError{reason: reason}
}

// Error returns the underlying failure reason and satisfies the error interface.
func (validationError *ValidationError) Error() string {
	return validationError.reason
}
