package model

// Paging Model
//
// The Paging holds information about the limit and making requests to the next page.
//
// swagger:bacnet_model Paging
type Paging struct {
	// The request url for the next page. Empty/Null when no next page is available.
	Next string `json:"next,omitempty"`
	// The amount of messages that got returned in the current request.
	Size int `json:"size"`
	// The ID of the last message returned in the current request. Use this as alternative to the next link.
	Since uint `json:"since"`
	// The limit of the messages for the current request.
	Limit int `json:"limit"`
}

// PagedMessages Model
//
// Wrapper for the paging and the messages
//
// swagger:bacnet_model PagedMessages
type PagedMessages struct {
	// The paging of the messages.
	Paging Paging `json:"paging"`
	// The messages.
	Messages []*MessageExternal `json:"messages"`
}
