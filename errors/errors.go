package errors

// TypeError is returned when the t key was either missing or invalid
type TypeError struct{}
func (t TypeError) Error() string {
	return "TypeError"
}

// DataError is returned when the t key was either missing or invalid
type DataError struct{}
func (t DataError) Error() string {
	return "DataError"
}

// IDError is returned when the t key was either missing or invalid
type IDError struct{}
func (t IDError) Error() string {
	return "IDError"
}