package repositories

import (
	"context"
	"fmt"
	"goozinshe/logger"
	"goozinshe/models"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type MoviesAdminRepository struct {
	db *pgxpool.Pool
}

func NewMoviesAdminRepository(conn *pgxpool.Pool) *MoviesAdminRepository {
	return &MoviesAdminRepository{db: conn}
}

func (r *MoviesAdminRepository) FindById(c context.Context, id int) (models.MovieAdminResponse, error) {
	sql :=
		`
SELECT 
        m.id,
        m.title,
        m.description,
        m.release_year,
        m.director,
        m.rating,
        m.is_watched,
        m.trailer_url,
        m.poster_url,
        g.id,
        g.title,
        g.poster_url,
        c.id,
        c.title,
        c.poster_url,
        a.id,
        a.age,
        a.poster_url,
        e.id, 
        e.series, 
        e.title,            
        e.description,
        e.release_year,
        e.director,      
        e.rating,         
        e.trailer_url
    FROM movies m
    JOIN movies_genres mg ON mg.movie_id = m.id
    JOIN genres g ON mg.genre_id = g.id
    JOIN movies_categories mc ON mc.movie_id = m.id
    JOIN categories c ON mc.categorie_id = c.id
    JOIN movies_ages ma ON ma.movie_id = m.id
	JOIN ages a ON ma.age_id = a.id
    left JOIN movies_allseries me ON me.movie_id = m.id
    left JOIN allseries e ON me.allserie_id = e.id
where m.id = $1
	`

	logger := logger.GetLogger()

	rows, err := r.db.Query(c, sql, id)
	defer rows.Close()
	if err != nil {
		logger.Error("Could not query database", zap.String("db_msg", err.Error()))
		return models.MovieAdminResponse{}, err
	}
	var movie *models.MovieAdminResponse

	categoriesMap := make(map[int]*models.Category, 0)
	category := make([]*models.Category, 0)

	genresMap := make(map[int]*models.Genre, 0)
	genre := make([]*models.Genre, 0)

	agesMap := make(map[int]*models.Age, 0)
	age := make([]*models.Age, 0)

	allseriesMap := make(map[int]*models.AllSeries, 0)
	allserie := make([]*models.AllSeries, 0)

	for rows.Next() {
		var m models.MovieAdminResponse
		var g models.Genre
		var c models.Category
		var a models.Age
		var e models.AllSeries

		err := rows.Scan(
			&m.Id,
			&m.Title,
			&m.Description,
			&m.ReleaseYear,
			&m.Director,
			&m.Rating,
			&m.IsWatched,
			&m.TrailerUrl,
			&m.PosterUrl,
			&g.Id,
			&g.Title,
			&g.PosterUrl,
			&c.Id,
			&c.Title,
			&c.PosterUrl,
			&a.Id,
			&a.Age,
			&a.PosterUrl,
			&e.Id,
			&e.Series,
			&e.Title,
			&e.Description,
			&e.ReleaseYear,
			&e.Director,
			&e.Rating,
			&e.TrailerUrl,
		)
		if err != nil {
			return models.MovieAdminResponse{}, err
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

		if e.Id != nil {
			if _, exists := allseriesMap[*e.Id]; !exists {
				allseriesMap[*e.Id] = &e
				allserie = append(allserie, &e)
			}
		}
	}

	err = rows.Err()
	if err != nil {
		return models.MovieAdminResponse{}, err
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

	var allseries []models.AllSeries
	for _, allserie := range allserie {
		allseries = append(allseries, *allserie)
	}

	movie.Category = categories
	movie.Genres = genres
	movie.Ages = ages
	movie.AllSeries = allseries

	return *movie, nil
}

func (r *MoviesAdminRepository) FindAll(c context.Context, filters models.MovieFilters) ([]models.Movie, error) {
	sql := ` 
    SELECT 
        m.id,
        m.title,
        m.description,
        m.release_year,
        m.director,
        m.rating,
        m.is_watched,
        m.trailer_url,
        m.poster_url,
        g.id,
        g.title,
        g.poster_url,
        c.id,
        c.title,
        c.poster_url,
        a.id,
        a.age,
        a.poster_url,
        e.id, 
        e.series, 
        e.title,            
        e.description,
        e.release_year,
        e.director,      
        e.rating,         
        e.trailer_url
    FROM movies m
    JOIN movies_genres mg ON mg.movie_id = m.id
    JOIN genres g ON mg.genre_id = g.id
    JOIN movies_categories mc ON mc.movie_id = m.id
    JOIN categories c ON mc.categorie_id = c.id
    JOIN movies_ages ma ON ma.movie_id = m.id
    JOIN ages a ON ma.age_id = a.id
    left JOIN movies_allseries me ON me.movie_id = m.id
    left JOIN allseries e ON me.allserie_id = e.id
    `

	params := pgx.NamedArgs{}

	if filters.SearchTerm != "" {
		// '%%%s%%' => '%поиск%'
		sql = fmt.Sprintf("%s and m.title ilike @s", sql)
		params["s"] = fmt.Sprintf("%%%s%%", filters.SearchTerm)
	}
	if filters.GenreId != "" {
		sql = fmt.Sprintf("%s and g.id = @genreId", sql)
		params["genreId"] = filters.GenreId
	}
	if filters.IsWatched != "" {
		isWatched, _ := strconv.ParseBool(filters.IsWatched)

		sql = fmt.Sprintf("%s and m.is_watched = @isWatched", sql)
		params["isWatched"] = isWatched
	}
	if filters.Sort != "" {
		identifier := pgx.Identifier{filters.Sort}
		sql = fmt.Sprintf("%s order by m.%s", sql, identifier.Sanitize())
	}

	l := logger.GetLogger()
	rows, err := r.db.Query(c, sql)
	if err != nil {
		l.Error(err.Error())
		return nil, err
	}

	movies := make([]*models.Movie, 0)
	moviesMap := make(map[int]*models.Movie)

	allseriesMap := make(map[int]*models.AllSeries, 0)

	for rows.Next() {
		var m models.Movie
		var g models.Genre
		var c models.Category
		var a models.Age
		var e models.AllSeries

		err := rows.Scan(
			&m.Id,
			&m.Title,
			&m.Description,
			&m.ReleaseYear,
			&m.Director,
			&m.Rating,
			&m.IsWatched,
			&m.TrailerUrl,
			&m.PosterUrl,
			&g.Id,
			&g.Title,
			&g.PosterUrl,
			&c.Id,
			&c.Title,
			&c.PosterUrl,
			&a.Id,
			&a.Age,
			&a.PosterUrl,
			&e.Id,
			&e.Series,
			&e.Title,
			&e.Description,
			&e.ReleaseYear,
			&e.Director,
			&e.Rating,
			&e.TrailerUrl,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := moviesMap[m.Id]; !exists {
			moviesMap[m.Id] = &m
			movies = append(movies, &m)
		}

		// Обработка жанров
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

		// Обработка категорий
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

		// Обработка возрастных категорий
		ageExists := false
		for _, existingAge := range moviesMap[m.Id].Ages {
			if existingAge.Id == a.Id {
				ageExists = true
				break
			}
		}
		if !ageExists {
			moviesMap[m.Id].Ages = append(moviesMap[m.Id].Ages, a)
		}
		// Обработка сериалов

		if e.Id != nil {
			if _, exists := allseriesMap[*e.Id]; !exists {
				allseriesMap[*e.Id] = &e
				moviesMap[m.Id].AllSeries = append(moviesMap[m.Id].AllSeries, e)
			}
		}
	}

	err = rows.Err()
	if err != nil {
		l.Error(err.Error())
		return nil, err
	}

	concreteMovies := make([]models.Movie, 0, len(movies))
	for _, v := range movies {
		concreteMovies = append(concreteMovies, *v)
	}

	return concreteMovies, nil
}

// allseriesExists := false
// for _, exisitingAllseries := range moviesMap[m.Id].AllSeries {
// 	if exisitingAllseries.Id == e.Id {
// 		allseriesExists = true
// 		break
// 	}
// }
// if !allseriesExists {
// 	moviesMap[m.Id].AllSeries = append(moviesMap[m.Id].AllSeries, e)
// }

func (r *MoviesAdminRepository) Create(c context.Context, movie models.MovieAdminResponse) (int, error) {
	l := logger.GetLogger()
	var id int

	tx, err := r.db.Begin(c)
	if err != nil {
		l.Error(err.Error())
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(c) // Если ошибка, откатываем транзакцию
		}
	}()

	row := tx.QueryRow(c,
		` 
    insert into movies(title, description, release_year, director, trailer_url, poster_url)
    values($1, $2, $3, $4, $5, $6)
    returning id
    `,
		movie.Title,
		movie.Description,
		movie.ReleaseYear,
		movie.Director,
		movie.TrailerUrl,
		movie.PosterUrl)

	err = row.Scan(&id)
	if err != nil {
		l.Error(err.Error())
		return 0, err
	}
	for _, genre := range movie.Genres {
		_, err = tx.Exec(c, "insert into movies_genres(movie_id, genre_id) values($1, $2)", id, genre.Id)
		if err != nil {
			l.Error(err.Error())
			return 0, err
		}
	}

	for _, category := range movie.Category {
		_, err = tx.Exec(c, "insert into movies_categories(movie_id, categorie_id) values($1, $2)", id, category.Id)
		if err != nil {
			l.Error(err.Error())
			return 0, err
		}
	}

	for _, age := range movie.Ages {
		_, err = tx.Exec(c, "insert into movies_ages(movie_id, age_id) values($1, $2)", id, age.Id)
		if err != nil {
			l.Error(err.Error())
			return 0, err
		}
	}

	for _, age := range movie.Ages {
		_, err = tx.Exec(c, "insert into movies_ages(movie_id, age_id) values($1, $2)", id, age.Id)
		if err != nil {
			l.Error(err.Error())
			return 0, err
		}
	}

	l.Info(fmt.Sprintf("проект %s добавлен успешно", movie.Title))

	err = tx.Commit(c)
	if err != nil {
		l.Error(err.Error())
		return 0, nil
	}

	return id, nil
}

func (r *MoviesAdminRepository) Update(c context.Context, id int, updatedMovie models.MovieAdminResponse) error {
	l := logger.GetLogger()
	tx, err := r.db.Begin(c)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(c) // Если ошибка, откатываем транзакцию
		}
	}()

	_, err = tx.Exec(
		c,
		`
        update movies
        set 
            title = $1,
            description = $2,
            release_year = $3,
            director = $4,
            trailer_url = $5,
            poster_url = $6
        where id = $7
        `,
		updatedMovie.Title,
		updatedMovie.Description,
		updatedMovie.ReleaseYear,
		updatedMovie.Director,
		updatedMovie.TrailerUrl,
		updatedMovie.PosterUrl,
		id)

	if err != nil {
		l.Error(err.Error())
		return err
	}

	_, err = tx.Exec(c, "DELETE FROM movies_genres WHERE movie_id = $1", id)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	_, err = tx.Exec(c, "DELETE FROM movies_categories WHERE movie_id = $1", id)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	_, err = tx.Exec(c, "DELETE FROM movies_ages WHERE movie_id = $1", id)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	for _, genre := range updatedMovie.Genres {
		_, err = tx.Exec(c, "insert into movies_genres(movie_id, genre_id) values($1, $2)", id, genre.Id)
		if err != nil {
			l.Error(err.Error())
			return err
		}
	}

	for _, category := range updatedMovie.Category {
		_, err = tx.Exec(c, "insert into movies_categories(movie_id, categorie_id) values($1, $2)", id, category.Id)
		if err != nil {
			l.Error(err.Error())
			return err
		}
	}

	for _, age := range updatedMovie.Ages {
		_, err = tx.Exec(c, "insert into movies_ages(movie_id, age_id) values($1, $2)", id, age.Id)
		if err != nil {
			l.Error(err.Error())
			return err
		}
	}

	err = tx.Commit(c)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	return nil
}

func (r *MoviesAdminRepository) Delete(c context.Context, id int) error {
	l := logger.GetLogger()
	tx, err := r.db.Begin(c)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	var movieTitle string
	row := r.db.QueryRow(c, "select title from movies where id = $1", id)
	err = row.Scan(&movieTitle)
	if err != nil {
		return err
	}

	l.Warn(fmt.Sprintf("Вы действительно хотите удалить %s?", movieTitle))

	_, err = tx.Exec(c, "DELETE FROM movies_genres WHERE movie_id = $1", id)
	if err != nil {
		l.Error(err.Error())
		tx.Rollback(c)
		return err
	}

	_, err = tx.Exec(c, "DELETE FROM movies_categories WHERE movie_id = $1", id)
	if err != nil {
		l.Error(err.Error())
		tx.Rollback(c)
		return err
	}

	_, err = tx.Exec(c, "DELETE FROM movies_ages WHERE movie_id = $1", id)
	if err != nil {
		l.Error(err.Error())
		tx.Rollback(c)
		return err
	}

	_, err = tx.Exec(c, "DELETE FROM movies_allseries WHERE movie_id = $1", id)
	if err != nil {
		l.Error(err.Error())
		tx.Rollback(c)
		return err
	}

	_, err = tx.Exec(c, "DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		l.Error(err.Error())
		tx.Rollback(c)
		return err
	}

	l.Info(fmt.Sprintf("Фильм %s удален:", movieTitle))

	// Фиксируем транзакцию
	err = tx.Commit(c)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	return nil
}

func (r *MoviesAdminRepository) SetWatched(c context.Context, id int, isWatched bool) error {
	_, err := r.db.Exec(c, "update movies set is_watched = $1 where id = $2", isWatched, id)
	if err != nil {
		l := logger.GetLogger()
		l.Error(err.Error())
		return err
	}

	return nil
}
