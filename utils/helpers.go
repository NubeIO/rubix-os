package utils

import (
	"fmt"
)

func truncateString(str string, num int) string {
	ret := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		ret = str[0:num] + ""
	}
	return ret
}
func NameIsNil() string {

	uuid := MakeTopicUUID("")
	return fmt.Sprintf("n_%s", truncateString(uuid, 8))

}
