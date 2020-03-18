package user

import (
	"fmt"

	"github.com/varrrro/recomovies/internal/database/rating"
)

// User present in the database.
type User struct {
	ID      int             `gorm:"primary_key"`
	Ratings []rating.Rating `gorm:"foreignkey:UserID"`
	Bias    float32
}

// New user instance.
func New(id int) *User {
	return &User{
		ID: id,
	}
}

// CalculateBias of the user based on his/her ratings.
func (u *User) CalculateBias() error {
	if u.Ratings == nil || len(u.Ratings) == 0 {
		return fmt.Errorf("The user has no ratings")
	}

	sum := 0
	for _, r := range u.Ratings {
		sum += r.Value
	}
	u.Bias = float32(sum) / float32(len(u.Ratings))
	return nil
}
