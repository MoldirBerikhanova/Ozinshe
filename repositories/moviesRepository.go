package repositories

import (
	"context"
	//"fmt"
	"goozinshe/logger"
	"goozinshe/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type MoviesRepository struct {
	db *pgxpool.Pool
}

func NewMoviesRepository(conn *pgxpool.Pool) *MoviesRepository {
	return &MoviesRepository{db: conn}
}

func (r *MoviesRepository) FindById(c context.Context, id int) (models.Movie, error) {
	sql :=
		`
select 
m.id,
m.title,
m.description,
m.release_year,
m.director,
m.rating,
m.is_watched,
m.trailer_url,
g.id,
g.title,
c.id,
c.title,
a.id,
a.age
from movies m
join movies_genres mg on mg.movie_id = m.id
join genres g on mg.genre_id  = g.id
join movies_categories mc on mc.movie_id = m.id
join categories c on mc.categorie_id  = c.id
join movies_ages ma on ma.movie_id =m.id
join ages a on ma.age_id = a.id
where m.id = $1
	`

	logger := logger.GetLogger()

	rows, err := r.db.Query(c, sql, id)
	defer rows.Close()
	if err != nil {
		logger.Error("Could not query database", zap.String("db_msg", err.Error()))
		return models.Movie{}, err
	}
	var movie *models.Movie

	categoriesMap := make(map[int]*models.Category, 0)
	category := make([]*models.Category, 0)

	genresMap := make(map[int]*models.Genre, 0)
	genre := make([]*models.Genre, 0)

	agesMap := make(map[int]*models.Age, 0)
	age := make([]*models.Age, 0)

	for rows.Next() {
		var m models.Movie
		var g models.Genre
		var c models.Category
		var a models.Age

		err := rows.Scan(
			&m.Id,
			&m.Title,
			&m.Description,
			&m.ReleaseYear,
			&m.Director,
			&m.Rating,
			&m.IsWatched,
			&m.TrailerUrl,
			&g.Id,
			&g.Title,
			&c.Id,
			&c.Title,
			&a.Id,
			&a.Age,
		)
		if err != nil {
			return models.Movie{}, err
		}

		if movie == nil {
			movie = &m
		}

		if _, exists := categoriesMap[c.Id]; !exists {
			categoriesMap[c.Id] = &c
			category = append(category, &c)
		}

		if _, exists := genresMap[g.Id]; !exists {
			genresMap[g.Id] = &g
			genre = append(genre, &g)
		}

		if _, exists := agesMap[a.Id]; !exists {
			agesMap[a.Id] = &a
			age = append(age, &a)
		}
		//movie.Genres = append(movie.Genres, g)
	}

	err = rows.Err()
	if err != nil {
		return models.Movie{}, err
	}

	var categories []models.Category
	for _, cat := range category {
		categories = append(categories, *cat)
	}

	var genres []models.Genre
	for _, gen := range genre {
		genres = append(genres, *gen)
	}

	var ages []models.Age
	for _, age := range age {
		ages = append(ages, *age)
	}

	movie.Category = categories
	movie.Genres = genres
	movie.Ages = ages

	return *movie, nil
}

func (r *MoviesRepository) FindAll(c context.Context) ([]models.Movie, error) {
	sql :=
		`
select 
m.id,
m.title,
m.description,
m.release_year,
m.director,
m.rating,
m.is_watched,
m.trailer_url,
g.id,
g.title,
c.id,
c.title,
a.id,
a.age,
r.id,
r.names_of_hero,
r.names_of_actors
from movies m
join movies_genres mg on mg.movie_id = m.id
join genres g on mg.genre_id  = g.id
join movies_categories mc on mc.movie_id = m.id
join categories c on mc.categorie_id  = c.id
join movies_ages ma on ma.movie_id =m.id
join ages a on ma.age_id = a.id
join movies_roles mr on mr.movie_id =m.id
join roles r on mr.role_id = r.id


	`

	rows, err := r.db.Query(c, sql)
	if err != nil {
		return nil, err
	}

	movies := make([]*models.Movie, 0)
	moviesMap := make(map[int]*models.Movie)

	for rows.Next() {
		var m models.Movie
		var g models.Genre
		var c models.Category
		var a models.Age
		var r models.Roles

		err := rows.Scan(
			&m.Id,
			&m.Title,
			&m.Description,
			&m.ReleaseYear,
			&m.Director,
			&m.Rating,
			&m.IsWatched,
			&m.TrailerUrl,
			&g.Id,
			&g.Title,
			&c.Id,
			&c.Title,
			&a.Id,
			&a.Age,
			&r.Id,
			&r.Names,
			&r.Actors,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := moviesMap[m.Id]; !exists {
			moviesMap[m.Id] = &m
			movies = append(movies, &m)

		}

		genreExists := false
		for _, existingGenres := range moviesMap[m.Id].Genres {
			if existingGenres.Id == g.Id {
				genreExists = true
				break
			}
		}

		if !genreExists {
			moviesMap[m.Id].Genres = append(moviesMap[m.Id].Genres, g)
		}

		categoryExists := false
		for _, existingCategory := range moviesMap[m.Id].Category {
			if existingCategory.Id == c.Id {
				categoryExists = true
				break
			}
		}

		if !categoryExists {
			moviesMap[m.Id].Category = append(moviesMap[m.Id].Category, c)
		}

		ageExists := false
		for _, exisitingAges := range moviesMap[m.Id].Ages {
			if exisitingAges.Id == a.Id {
				ageExists = true
				break
			}
		}
		if !ageExists {
			moviesMap[m.Id].Ages = append(moviesMap[m.Id].Ages, a)
		}

		roleExists := false
		for _, exisitingRoles := range moviesMap[m.Id].Roles {
			if exisitingRoles.Id == r.Id {
				roleExists = true
				break
			}
		}
		if !roleExists {
			moviesMap[m.Id].Roles = append(moviesMap[m.Id].Roles, r)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	concreteMovies := make([]models.Movie, 0, len(movies))
	for _, v := range movies {
		concreteMovies = append(concreteMovies, *v)
	}

	return concreteMovies, nil
}

func (r *MoviesRepository) Create(c context.Context, movie models.Movie) (int, error) {
	var id int

	tx, err := r.db.Begin(c)
	if err != nil {
		return 0, nil
	}
	row := tx.QueryRow(c,
		`
insert into movies(title, description, release_year, director, trailer_url)
values($1, $2, $3, $4, $5)
returning id
	`,
		movie.Title,
		movie.Description,
		movie.ReleaseYear,
		movie.Director,
		movie.TrailerUrl)

	err = row.Scan(&id)
	if err != nil {
		return 0, nil
	}
	//Вставка жанров для фильма в таблицу movies_genres
	for _, genre := range movie.Genres {
		_, err = tx.Exec(c, "insert into movies_genres(movie_id, genre_id) values($1, $2)", id, genre.Id)
		if err != nil {
			return 0, err
		}
	}

	for _, category := range movie.Category {
		_, err = tx.Exec(c, "insert into movies_categories(movie_id, categorie_id) values($1, $2)", id, category.Id)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit(c)
	if err != nil {
		return 0, nil
	}

	return id, nil
}

func (r *MoviesRepository) Update(c context.Context, id int, updatedMovie models.Movie) error {
	tx, err := r.db.Begin(c)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		c,
		`
update movies
set 
title = $1,
description = $2,
release_year = $3,
director = $4,
trailer_url = $5
where id = $6
	`,
		updatedMovie.Title,
		updatedMovie.Description,
		updatedMovie.ReleaseYear,
		updatedMovie.Director,
		updatedMovie.TrailerUrl,
		id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(c, "delete from movies_genres where movie_id = $1", id)
	if err != nil {
		return err
	}
	for _, genre := range updatedMovie.Genres {
		_, err = r.db.Exec(c, "insert into movies_genres(movie_id, genre_id) values($1, $2)", id, genre.Id)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(c)
	if err != nil {
		return err
	}

	return nil
}

func (r *MoviesRepository) Delete(c context.Context, id int) error {
	tx, err := r.db.Begin(c)
	if err != nil {
		return err
	}

	_, err = tx.Exec(c, "delete from movies_genres where movie_id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(c, "delete from movies where id = $1", id)
	if err != nil {
		return err
	}

	err = tx.Commit(c)
	if err != nil {
		return err
	}

	return nil
}
