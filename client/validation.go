package client

// Validatable is an interface for request types that support validation.
// Types implementing this interface will have Validate() called automatically
// before the request is sent, unless validation is disabled via WithNoValidation().
type Validatable interface {
	Validate() error
}
