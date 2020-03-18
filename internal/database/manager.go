package database

import (
	"github.com/jinzhu/gorm"

	"github.com/varrrro/recomovies/internal/database/movie"
	"github.com/varrrro/recomovies/internal/database/rating"
	"github.com/varrrro/recomovies/internal/database/user"
)

// Manager that controls database access.
type Manager struct {
	db *gorm.DB
}

// NewManager instance.
func NewManager(db *gorm.DB) *Manager {
	return &Manager{db}
}

// UpdateSchemas for the database.
func (mg *Manager) UpdateSchemas() {
	mg.db.AutoMigrate(&movie.Movie{}, &user.User{}, &rating.Rating{})
}

// AddMovie to database.
func (mg *Manager) AddMovie(m *movie.Movie) {
	mg.db.Create(m)
}

// AddUser to database.
func (mg *Manager) AddUser(u *user.User) {
	mg.db.Create(u)
}

// AddRating to database. Returns an error if either the movie
// or the user don't exist.
func (mg *Manager) AddRating(r *rating.Rating) error {
	var m movie.Movie
	mg.db.First(&m, "id = ?", r.MovieID)
	if m.ID != r.MovieID {
		return &NotFoundError{"movie", r.MovieID}
	}

	var u user.User
	mg.db.First(&u, "id = ?", r.UserID)
	if u.ID != r.UserID {
		return &NotFoundError{"user", r.UserID}
	}

	mg.db.Model(&m).Association("Ratings").Append(r)
	return nil
}

// FetchMovie with the given ID.
func (mg *Manager) FetchMovie(id int) (*movie.Movie, error) {
	var m movie.Movie
	mg.db.Preload("Ratings").First(&m, "id = ?", id)
	if m.ID != id {
		return &movie.Movie{}, &NotFoundError{"movie", id}
	}
	return &m, nil
}

// FetchUser with the given ID.
func (mg *Manager) FetchUser(id int) (*user.User, error) {
	var u user.User
	mg.db.Preload("Ratings").First(&u, "id = ?", id)
	if u.ID != id {
		return &user.User{}, &NotFoundError{"user", id}
	}
	return &u, nil
}

// FetchAllUsers present in the database.
func (mg *Manager) FetchAllUsers() []user.User {
	var u []user.User
	mg.db.Preload("Ratings").Find(&u)
	return u
}

// GetMovieCount of the database.
func (mg *Manager) GetMovieCount() int {
	var count int
	mg.db.Model(&movie.Movie{}).Count(&count)
	return count
}
