package session

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/varrrro/recomovies/internal/database/movie"
	"github.com/varrrro/recomovies/internal/database/user"
)

var (
	current = make(map[string]*Session)
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// GetSession from memory.
func GetSession(id string) (*Session, error) {
	s, ok := current[id]
	if !ok {
		return nil, fmt.Errorf("No session with ID %s", id)
	}
	return s, nil
}

// Session in progress.
type Session struct {
	ID              string
	Movies          []int
	Ratings         map[int]int
	Bias            float64
	Neighborhood    []*user.User
	Recommendations []*movie.Movie
}

// New session instance.
func New() *Session {
	s := &Session{
		ID:      randomID(6),
		Movies:  make([]int, 20),
		Ratings: make(map[int]int),
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
