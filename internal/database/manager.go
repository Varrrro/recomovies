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
	if err := mg.db.First(&m, "id = ?", r.MovieID).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	var u user.User
	if err := mg.db.First(&u, "id = ?", r.UserID).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	if err := mg.db.Model(&m).Association("Ratings").Append(r).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}
	return nil
}

// FetchMovie with the given ID.
func (mg *Manager) FetchMovie(id int) (*movie.Movie, error) {
	var m movie.Movie
	if err := mg.db.Preload("Ratings").First(&m, "id = ?", id).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &m, nil
}

// FetchUser with the given ID.
func (mg *Manager) FetchUser(id int) (*user.User, error) {
	var u user.User
	if err := mg.db.Preload("Ratings").First(&u, "id = ?", id).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &u, nil
}

// FetchAllUsers present in the database.
func (mg *Manager) FetchAllUsers() []*user.User {
	var u []*user.User
	mg.db.Preload("Ratings").Find(&u)
	return u
}

// FetchUserRating for a given movie, if it exists.
func (mg *Manager) FetchUserRating(userID, movieID int) (*rating.Rating, error) {
	var r rating.Rating
	if err := mg.db.Where(&rating.Rating{UserID: userID, MovieID: movieID}).First(&r).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &r, nil
}

// FetchUnseenMoviesRatings by users in a given group.
func (mg *Manager) FetchUnseenMoviesRatings(movieIDs, userIDs []int) []*rating.Rating {
	var r []*rating.Rating
	mg.db.Where("user_id IN (?)", userIDs).Not("movie_id", movieIDs).Find(&r)
	return r
}

// GetMovieCount of the database.
func (mg *Manager) GetMovieCount() int {
	var count int
	mg.db.Model(&movie.Movie{}).Count(&count)
	return count
}
