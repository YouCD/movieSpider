package model

import (
	"fmt"
	"testing"
)

func TestMovieDB_GetFeedVideoMovieByNames(t *testing.T) {
	movieDB := NewMovieDB()

	got, err := movieDB.findMovie([]string{"The.Fogsa"}, `name = ? and magnet!="" and type="movie"`)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(got)
}
