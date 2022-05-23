package connection

type Connection int64

const (
	Connected Connection = iota
	Broken
)

func (s Connection) String() string {
	switch s {
	case Connected:
		return "Connected"
	case Broken:
		return "Broken"
	}
	return "Broken"
}
