package interfaces

import "time"

type Snapshots struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
}
type CreateSnapshot struct {
	Description string `json:"description"`
}

type RestoreSnapshot struct {
	File        string `json:"file"`
	Description string `json:"description"`
}
