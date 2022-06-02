package nerrors

type ErrConflict struct {
	msg string
}

func (e *ErrConflict) Error() string { return e.msg }

func NewErrConflict(text string) error {
	return &ErrConflict{msg: text}
}
