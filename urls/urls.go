package urls

import "fmt"

var (
	producerURL = "/api/producers"
)

func ProducerURL() string {
	return producerURL
}

func ProducerSingularURL(uuid string) string {
	return fmt.Sprintf("%s/%s", producerURL, uuid)
}
