package lobby

type Canvas struct {
	Name      string
	Positions map[string]*Position
}

type Position struct {
	X int
	Y int
}

func NewCanvas(name string) *Canvas {
	return &Canvas{
		Name:      name,
		Positions: make(map[string]*Position),
	}
}
