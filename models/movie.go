package models

type MovieFilters struct {
	SearchTerm string
	GenreId    string
	IsWatched  string
	Sort       string
}

type Movie struct {
	Id          int         `form:"id"`
	Title       string      `form:"title"`
	Description string      `form:"description"`
	ReleaseYear int         `form:"release_year"`
	Director    string      `form:"director"`
	Rating      int         `form:"rating"`
	IsWatched   bool        `form:"is_watched"`
	TrailerUrl  string      `form:"trailer_url"`
	PosterUrl   string      `form:"poster_url"`
	Genres      []Genre     `form:"genres"`
	Category    []Category  `form:"categories"`
	Ages        []Age       `form:"ages"`
	AllSeries   []AllSeries `form:"allseries"`
}

type MovieAdminResponse struct {
	Id          int         `form:"id"`
	Title       string      `form:"title"`
	Description string      `form:"description"`
	ReleaseYear int         `form:"release_year"`
	Director    string      `form:"director"`
	Rating      int         `form:"rating"`
	IsWatched   bool        `form:"is_watched"`
	TrailerUrl  string      `form:"trailer_url"`
	PosterUrl   string      `form:"poster_url"`
	Genres      []Genre     `form:"genres"`
	Category    []Category  `form:"categories"`
	Ages        []Age       `form:"ages"`
	AllSeries   []AllSeries `form:"allseries"`
}
