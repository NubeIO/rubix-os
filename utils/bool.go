package utils

func NewTrue() *bool {
	b := true
	return &b
}

func NewFalse() *bool {
	b := false
	return &b
}

func BoolIsNil(b *bool) bool {
	if b == nil {
		return false
	} else {
		return *b
	}
}

func IsTrue(b *bool) bool {
	if b == nil {
		return false
	} else {
		return *b
	}
}
