package errors

type wrappedError interface{
	error
	SetError(err error)
	HasError() bool
}

type WrappedError struct {
	err error
}

func (w *WrappedError) Error() string {
	if w.err != nil {
		return w.err.Error()
	}

	return ""
}

func (w *WrappedError) SetError(err error) {
	w.err = err
}

func (w *WrappedError) HasError() bool {
	return w.err != nil
}