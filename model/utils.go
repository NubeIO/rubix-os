package model

func IsFNCreator(t interface{}) bool {
	// when it has pointer of value
	_, ok := t.(*FlowNetwork)
	if ok {
		return true
	}
	// when it has value
	switch t.(type) {
	case FlowNetwork:
		return true
	default:
		return false
	}
}
