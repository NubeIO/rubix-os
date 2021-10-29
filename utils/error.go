package utils

func NewError(err error) *error {
	e := err
	return &e
}
