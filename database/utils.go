package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
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

func typeIsNil(t string, use string) string {
	if t == "" {
		return use
	}
	return t
}

func nameIsNil(name string) string {
	if name == "" {
		uuid := utils.MakeTopicUUID("")
		return fmt.Sprintf("n_%s", truncateString(uuid, 8))
	}
	return name
}

func pluginIsNil(name string) string {
	if name == "" {
		return "system"
	}
	return name
}