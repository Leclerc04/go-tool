package errs

func RunAndRecover(action func() error) (errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = WrapPanicValue(r)
		}
	}()
	errOut = action()
	return errOut
}
