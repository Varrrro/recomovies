package handler

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/varrrro/recomovies/internal/database/movie"
	"github.com/varrrro/recomovies/internal/database/user"
	"github.com/varrrro/recomovies/internal/session"
)

const (
	k = 10
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
	s.Bias = float64(sum) / float64(len(s.Ratings))
}

func (h *Handler) defineNeighbors(s *session.Session) {
	users := h.mg.FetchAllUsers()
	temp := make([]*user.User, 0)
	for _, u := range users {
		matches := make(map[int]int)
		for _, id := range s.Movies {
			if r, err := h.mg.FetchUserRating(u.ID, id); err == nil {
				matches[id] = r.Value
			}
		}

		// Calculate similarity
		if len(matches) >= n {
			sumProduct, sumSessionSquare, sumUserSquare := 0.0, 0.0, 0.0
			for id, r := range matches {
				sessionValue := float64(s.Ratings[id]) - s.Bias
				userValue := float64(r) - u.Bias

				sumProduct += sessionValue * userValue
				sumSessionSquare += math.Pow(sessionValue, 2)
				sumUserSquare += math.Pow(userValue, 2)
			}
			u.Similarity = sumProduct / (math.Sqrt(sumSessionSquare) * math.Sqrt(sumUserSquare))

			if temp == nil || len(temp) == 0 {
				temp = append(temp, u)
			} else {
				i := sort.Search(len(temp), func(i int) bool {
					if u.Similarity <= temp[i].Similarity {
						return true
					}
					return false
				})
				temp = append(temp, nil)
				copy(temp[i+1:], temp[i:])
				temp[i] = u
			}
		}
	}

	//
	if len(temp) <= k {
		s.Neighborhood = make([]*user.User, len(temp))
		copy(s.Neighborhood, temp)
	} else {
		s.Neighborhood = make([]*user.User, k)
		copy(s.Neighborhood, temp[len(temp)-k:])
	}
}

func (h *Handler) obtainRecommendations(s *session.Session) {
	sum := 0.0
	neighbors := make([]int, 0)
	for _, u := range s.Neighborhood {
		sum += math.Abs(u.Similarity)
		neighbors = append(neighbors, u.ID)
	}
	c := 1.0 / sum

	// Calculate predicted ratings for unseen movies
	temp := make(map[*movie.Movie]float64)
	ratings := h.mg.FetchUnseenMoviesRatings(s.Movies, neighbors)
	for _, r := range ratings {
		sum := 0.0
		for _, u := range s.Neighborhood {
			if ur, err := h.mg.FetchUserRating(u.ID, r.MovieID); err == nil {
				sum += u.Similarity * (float64(ur.Value) - u.Bias)
			}
		}

		if m, err := h.mg.FetchMovie(r.MovieID); err == nil {
			pr := s.Bias + c*sum
			temp[m] = pr
		}
	}

	// Find max value
	max := 0.0
	for _, pr := range temp {
		if pr > max {
			max = pr
		}
	}

	// Rescale predicted ratings to [1, 5] interval
	for m, pr := range temp {
		scaledpr := ((pr-1)/(max-1))*4 + 1
		if scaledpr >= 4.0 {
			s.Recommendations = append(s.Recommendations, m)
		}
	}
}
