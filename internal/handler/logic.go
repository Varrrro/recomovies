package handler

import (
	"math/rand"
	"time"

	"github.com/varrrro/recomovies/internal/session"
)

const (
	k = 20
	n = 5
)

func (h *Handler) initMovies(s *session.Session) {
	rand.Seed(time.Now().UnixNano())
	n := h.mg.GetMovieCount()
	for i := 0; i < 20; i++ {
		s.Movies[i] = rand.Intn(n) + 1
	}
}

func (h *Handler) calculateBias(s *session.Session) {
	sum := 0
	for _, r := range s.Ratings {
		sum += r
	}
	s.Bias = float32(sum) / float32(len(s.Ratings))
}

func (h *Handler) defineNeighbors(s *session.Session) {

}

func (h *Handler) obtainRecommendations(s *session.Session) {

}
