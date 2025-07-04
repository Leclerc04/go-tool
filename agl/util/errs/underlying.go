package errs

// Unwrap unwraps the error, and return the underlying error.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	if e2, ok := err.(*Error); ok {
		if e2.Err == nil {
			return e2
		}
		return Unwrap(e2.Err)
	}
	return err
}
