package utils

import "path"

func ParseUUID(url string) (string, string) {
	parentURL, uuid := path.Split(url)
	return path.Clean(parentURL), uuid
}
