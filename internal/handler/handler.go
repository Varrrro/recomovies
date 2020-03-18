package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/varrrro/recomovies/internal/database"
	"github.com/varrrro/recomovies/internal/database/movie"
	"github.com/varrrro/recomovies/internal/session"
)

// Handler of http requests.
type Handler struct {
	mg *database.Manager
}

// New handler instance.
func New(mg *database.Manager) *Handler {
	return &Handler{mg}
}

// GetIndex handler.
func (h *Handler) GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// PostIndex handler.
func (h *Handler) PostIndex(c *gin.Context) {
	s := session.New()
	h.initMovies(s)

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/rating/1?session_id=%s", s.ID))
}

// GetRating handler.
func (h *Handler) GetRating(c *gin.Context) {
	s, _ := session.GetSession(c.Query("session_id"))
	n, _ := strconv.Atoi(c.Param("n"))
	m, _ := h.mg.FetchMovie(s.Movies[n-1])

	c.HTML(http.StatusOK, "rate.html", gin.H{
		"session_id": s.ID,
		"title":      m.Title,
		"n":          n,
	})
}

// PostRating handler.
func (h *Handler) PostRating(c *gin.Context) {
	s, _ := session.GetSession(c.Query("session_id"))
	n, _ := strconv.Atoi(c.Param("n"))
	r, _ := strconv.Atoi(c.PostForm("rating"))

	s.Ratings[s.Movies[n-1]] = r

	var redirectURL string
	if n < 20 {
		redirectURL = fmt.Sprintf("/rating/%d?session_id=%s", n+1, s.ID)
	} else {
		redirectURL = fmt.Sprintf("/results?session_id=%s", s.ID)

		// Process ratings and get recommendations
		h.calculateBias(s)
		h.defineNeighbors(s)
		h.obtainRecommendations(s)
	}
	c.Redirect(http.StatusMovedPermanently, redirectURL)
}

// GetResults handler.
func (h *Handler) GetResults(c *gin.Context) {
	s, _ := session.GetSession(c.Query("session_id"))
	recs := make([]*movie.Movie, len(s.Recommendations))
	for i := 0; i < len(s.Recommendations); i++ {
		recs[i], _ = h.mg.FetchMovie(s.Recommendations[i])
	}

	c.HTML(http.StatusOK, "results.html", recs)
}
