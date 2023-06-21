package schema

type PollPriority struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Poll Priority"`
	Options  []string `json:"enum" default:"[\"high\",\"normal\",\"low\"]"`
	EnumName []string `json:"enumNames" default:"[\"high\",\"normal\",\"low\"]"`
	Default  string   `json:"default" default:"normal"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type PollRate struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Poll Rate"`
	Options  []string `json:"enum" default:"[\"fast\",\"normal\",\"slow\"]"`
	EnumName []string `json:"enumNames" default:"[\"fast\",\"normal\",\"slow\"]"`
	Default  string   `json:"default" default:"normal"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type FastPollRate struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Fast Poll Rate (seconds)"`
	Default  int    `json:"default" default:"1"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type NormalPollRate struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Normal Poll Rate (seconds)"`
	Default  int    `json:"default" default:"20"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type SlowPollRate struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Slow Poll Rate (seconds)"`
	Default  int    `json:"default" default:"120"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type MaxPollRate struct {
	Type     string  `json:"type" default:"number"`
	Title    string  `json:"title" default:"Max Poll Rate (seconds)"`
	Default  float64 `json:"default" default:"0.1"`
	ReadOnly bool    `json:"readOnly" default:"false"`
}

type ZeroMode struct {
	Type     string `json:"type" default:"boolean"`
	Title    string `json:"title" default:"Zero Mode"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type TimeoutSecs struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Timeout (seconds)"`
	Default  int    `json:"default" default:"2"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}
