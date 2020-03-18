package rating

// Rating of a movie given by an user.
type Rating struct {
	UserID  int
	MovieID int
	Value   int
}

// New rating instance.
func New(userID, movieID, value int) *Rating {
	return &Rating{
		UserID:  userID,
		MovieID: movieID,
		Value:   value,
	}
}
