package lobby

import (
	"math/rand"
	"sync"
	"time"
)

type Canvas struct {
	Name      string
	Positions map[string]*Position
	mu        sync.RWMutex
}

type Position struct {
	Username string
	X        float32
	Y        float32
}

const (
	minX = 1
	maxX = 399
	minY = 1
	maxY = 299
)

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

// удалить чела из коорд мапы
func (ca *Canvas) RemoveUser(name string) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	delete(ca.Positions, name)
}

func (ca *Canvas) AddUser(name string, x, y float32) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	ca.Positions[name] = NewPositionWithBounds(name, x, y)
}

func (ca *Canvas) ChangeUserCoords(name string, x, y float32) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	ca.Positions[name] = NewPositionWithBounds(name, x, y)
}

func NewPositionWithBounds(name string, x, y float32) *Position {

	// Ограничение значений X и Y в заданных диапазонах
	if x < minX {
		x = minX
	} else if x > maxX {
		x = maxX
	}

	if y < minY {
		y = minY
	} else if y > maxY {
		y = maxY
	}

	return &Position{
		Username: name,
		X:        x,
		Y:        y,
	}
}

func NewPositionRandomCoords(name string) *Position {
	rand.Seed(time.Now().UnixNano())

	randomX := float32(rand.Intn(maxX-minX+1) + minX)

	randomY := float32(rand.Intn(maxY-minY+1) + minY)
	return &Position{
		Username: name,
		X:        randomX,
		Y:        randomY,
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
