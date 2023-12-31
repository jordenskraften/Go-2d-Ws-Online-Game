package lobby

import "sync"

type Canvas struct {
	Name      string
	Positions map[string]*Position
	mu        sync.RWMutex
}

type Position struct {
	Name string
	X    float32
	Y    float32
}

func NewCanvas(name string) *Canvas {
	return &Canvas{
		Name:      name,
		Positions: make(map[string]*Position),
	}
}

//есть ли чел в канвасе?

// Check if a user is in the canvas
func (ca *Canvas) IsUserInCanvas(name string) bool {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	_, exists := ca.Positions[name]
	return exists
}

// добавить чела в мап и назначить стартовые коорды
func (ca *Canvas) AddUser(name string, x float32, y float32) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	ca.Positions[name] = &Position{
		X: x,
		Y: y,
	}
}

//удалить чела из коорд мапы
func (ca *Canvas) RemoveUser(name string) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	delete(ca.Positions, name)
}

//поменять челу коорды
func (ca *Canvas) ChangeUserCoords(name string, x float32, y float32) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	ca.Positions[name] = &Position{
		X: x,
		Y: y,
	}
}

// выслать инфо о канвасе
func (ca *Canvas) GetCanvasInfo() []*Position {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	info := make([]*Position, 0, len(ca.Positions))

	for _, val := range ca.Positions {
		info = append(info, val)
	}

	return info
}
