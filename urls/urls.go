package urls

import "fmt"

func ProducerURL() string {
	return "/api/producers"
}

func ProducerURLWithStream(streamUUID string) string {
	return fmt.Sprintf("%s?stream_uuid=%s", ProducerURL(), streamUUID)
}

func ProducerSingularURL(uuid string) string {
	return fmt.Sprintf("%s/%s", ProducerURL(), uuid)
}
