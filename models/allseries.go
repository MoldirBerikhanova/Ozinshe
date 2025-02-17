package models

type AllSeries struct {
	Id          *int    `form:"id"` // Указатель на int
	Series      *int    `form:"series"`
	Title       *string `form:"title"`
	Description *string `form:"description"`
	ReleaseYear *int    `form:"release_year"`
	Director    *string `form:"director"`
	Rating      *int    `form:"rating"`
	TrailerUrl  *string `form:"trailer_url"`
}
