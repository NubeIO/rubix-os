package boolean

func New(value bool) *bool {
	b := value
	return &b
}

func NewTrue() *bool {
	b := true
	return &b
}

func NewFalse() *bool {
	b := false
	return &b
}

func IsTrue(b *bool) bool {
	if b == nil {
		return false
	} else {
		return *b
	}
}

func IsFalse(b *bool) bool {
	return !IsTrue(b)
}

func IsNil(b *bool) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}
