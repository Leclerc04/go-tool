package errs

// Ignore ignores the err. You typically won't need this.
// This should only be used if the err
// is properly handled in a different way.
func Ignore(err error) {}
