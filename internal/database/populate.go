package database

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/varrrro/recomovies/internal/database/movie"
	"github.com/varrrro/recomovies/internal/database/rating"
	"github.com/varrrro/recomovies/internal/database/user"
)

const (
	moviesPath  = "data/u.item"
	ratingsPath = "data/u.data"
)

// Populate database with data from files.
func Populate(mg *Manager) error {
	// Read movies
	moviesFile, err := os.Open(moviesPath)
	if err != nil {
		return err
	}
	defer moviesFile.Close()

	scanner := bufio.NewScanner(moviesFile)
	scanner.Split(bufio.ScanLines)

	log.Println("Reading movies")
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "|")
		id, _ := strconv.Atoi(line[0])

		m := &movie.Movie{
			ID:    id,
			Title: line[1],
		}

		mg.AddMovie(m)
	}

	// Read ratings
	ratingsFile, err := os.Open(ratingsPath)
	if err != nil {
		return err
	}
	defer ratingsFile.Close()

	scanner = bufio.NewScanner(ratingsFile)
	scanner.Split(bufio.ScanLines)

	log.Println("Reading users and ratings")
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		userid, _ := strconv.Atoi(line[0])
		movieid, _ := strconv.Atoi(line[1])
		value, _ := strconv.Atoi(line[2])

		if u, err := mg.FetchUser(userid); err != nil {
			u = &user.User{
				ID: userid,
			}
			mg.AddUser(u)
		}

		r := &rating.Rating{
			UserID:  userid,
			MovieID: movieid,
			Value:   value,
		}
		mg.AddRating(r)
	}

	calculateBiases(mg)

	return nil
}

func calculateBiases(mg *Manager) {
	users := mg.FetchAllUsers()
	for _, u := range users {
		if err := u.CalculateBias(); err != nil {
			log.Fatalf("Can't calculate bias for user %d , error: %s\n", u.ID, err.Error())
		}
	}
}
