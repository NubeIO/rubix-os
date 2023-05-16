package interfaces

type Message struct {
	Message string `json:"message"`
}

type FoundMessage struct {
	Found bool `json:"found"`
}
