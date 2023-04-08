package api

// The DBDatabase interface for encapsulating database access.
type DBDatabase interface {
	SyncTopics() // sync all the topics into the event bus
}
type DatabaseAPI struct {
	DB DBDatabase
}

func (a *DatabaseAPI) SyncTopics() {
	a.DB.SyncTopics()
}
