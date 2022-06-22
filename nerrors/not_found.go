package nerrors

type NotFound struct {
	msg string
}

func (e *NotFound) Error() string { return e.msg }

func NewNotFound(text string) error {
	return &NotFound{msg: text}
}
