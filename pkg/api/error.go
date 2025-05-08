package api

// HandleError wraps an error with an associated HTTP status code for standardized API responses.
type HandleError struct {
	Err    error // Underlying error
	Status int   // Corresponding HTTP status code
}
