package interfaces

type StreamLog struct {
	UUID           string   `json:"uuid"`
	Service        string   `json:"service" binding:"required"`
	Duration       int      `json:"duration" binding:"required"`
	KeyWordsFilter []string `json:"key_words_filter"` // example: mqtt, connected
	Message        []string `json:"message"`
}
