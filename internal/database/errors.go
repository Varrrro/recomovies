package database

import "fmt"

// NotFoundError used when querying the database.
type NotFoundError struct {
	item string
	id   int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("No %s found with ID %d", e.item, e.id)
}
