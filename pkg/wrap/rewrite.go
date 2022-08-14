package wrap

func Rewrite(err error, msg string) error {
	return &wrapErr{inner: err, msg: msg}
}

type wrapErr struct {
	inner error
	msg   string
}

func (w *wrapErr) Unwrap() error {
	return w.inner
}

func (w *wrapErr) Error() string {
	return w.msg
}
