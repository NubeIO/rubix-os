package schema

type Defaults struct {
	UUID        UUID        `json:"uuid"`
	Name        Name        `json:"name"`
	Description Description `json:"description"`
	Enable      Enable      `json:"enable"`
}

func GetDefaults() *Defaults {
	m := &Defaults{}
	Set(m)
	return m
}
