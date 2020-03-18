package movie

import "github.com/varrrro/recomovies/internal/database/rating"

// Movie present in the database.
type Movie struct {
	ID      int             `gorm:"primary_key"`
	Title   string          `gorm:"not null"`
	Ratings []rating.Rating `gorm:"foreignkey:MovieID"`
}

// New movie instance.
func New(id int, title string) *Movie {
	return &Movie{
		ID:    id,
		Title: title,
	}
}
