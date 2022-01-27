package utils

func FirstNotNilFloat(values ...*float64) *float64 {
	for _, n := range values {
		if n != nil {
			return n
		}
	}
	return nil
}
