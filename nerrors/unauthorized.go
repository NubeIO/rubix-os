package nerrors

type ErrUnauthorized struct {
	msg string
}

func (e *ErrUnauthorized) Error() string { return e.msg }

func NewErrUnauthorized(text string) error {
	return &ErrUnauthorized{msg: text}
}
