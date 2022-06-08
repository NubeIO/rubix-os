package interfaces

type Token struct {
	Name    string `json:"name" binding:"required"`
	Blocked *bool  `json:"blocked"`
}
