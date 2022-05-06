package boolean

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

func IsNil(b *bool) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}
