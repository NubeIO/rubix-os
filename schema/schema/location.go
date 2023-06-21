package schema

type Address struct {
	Type  string `json:"type" default:"string"`
	Title string `json:"title" default:"Address"`
	Min   int    `json:"minLength" default:"0"`
	Max   int    `json:"maxLength" default:"100"`
}

type City struct {
	Type  string `json:"type" default:"string"`
	Title string `json:"title" default:"City"`
	Min   int    `json:"minLength" default:"0"`
	Max   int    `json:"maxLength" default:"100"`
}

type State struct {
	Type  string `json:"type" default:"string"`
	Title string `json:"title" default:"State"`
	Min   int    `json:"minLength" default:"0"`
	Max   int    `json:"maxLength" default:"100"`
}

type Zip struct {
	Type  string `json:"type" default:"string"`
	Title string `json:"title" default:"Zip"`
	Min   int    `json:"minLength" default:"0"`
	Max   int    `json:"maxLength" default:"20"`
}

type Country struct {
	Type  string `json:"type" default:"string"`
	Title string `json:"title" default:"Country"`
	Min   int    `json:"minLength" default:"0"`
	Max   int    `json:"maxLength" default:"100"`
}

type Lat struct {
	Type  string `json:"type" default:"string"`
	Title string `json:"title" default:"Lat"`
	Min   int    `json:"minLength" default:"0"`
	Max   int    `json:"maxLength" default:"20"`
}

type Lon struct {
	Type  string `json:"type" default:"string"`
	Title string `json:"title" default:"Lon"`
	Min   int    `json:"minLength" default:"0"`
	Max   int    `json:"maxLength" default:"20"`
}

type Timezone struct {
	Type    string   `json:"type" default:"string"`
	Title   string   `json:"title" default:"Timezone"`
	Options []string `json:"enum"`
	Default string   `json:"default" default:"Africa/Abidjan"`
}
