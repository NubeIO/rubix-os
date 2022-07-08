package interfaces

type TokenCreate struct {
	Name    string `json:"name" binding:"required"`
	Blocked *bool  `json:"blocked" binding:"required"`
}

type TokenBlock struct {
	Blocked *bool `json:"blocked" binding:"required"`
}
