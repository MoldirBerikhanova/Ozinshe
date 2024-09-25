package models

type Movie struct {
	Id          int
	Title       string
	Description string
	ReleaseYear int
	Director    string
	Rating      int
	IsWatched   bool
	TrailerUrl  string
	PosterUrl   string  // not implemented
	Genres      []Genre // not implemented
}
