package interfaces

type RebootJob struct {
	Tag        string `json:"tag"`
	Expression string `json:"expression"`
}

type RestartJob struct {
	Unit       string `json:"unit"`
	Expression string `json:"expression"`
}
