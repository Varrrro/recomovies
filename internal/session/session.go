package session

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	current = make(map[string]*Session)
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// GetSession from memory.
func GetSession(id string) (*Session, error) {
	s, ok := current["id"]
	if !ok {
		return &Session{}, fmt.Errorf("No session with ID %s", id)
	}
	return s, nil
}

// Session in progress.
type Session struct {
	ID              string
	Movies          []int
	Ratings         map[int]int
	Bias            float32
	Neighborhood    map[int]float32
	Recommendations []int
}

// New session instance.
func New() *Session {
	s := &Session{
		ID:           randomID(6),
		Movies:       make([]int, 20),
		Ratings:      make(map[int]int),
		Neighborhood: make(map[int]float32),
	}
	current[s.ID] = s
	return s
}

func randomID(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
