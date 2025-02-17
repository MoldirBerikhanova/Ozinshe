package repositories

import (
	"context"
	"fmt"
	"goozinshe/logger"
	"goozinshe/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AllSeriesRepository struct {
	db *pgxpool.Pool
}

func NewAllSeriesRepository(conn *pgxpool.Pool) *AllSeriesRepository {
	return &AllSeriesRepository{db: conn}
}

func (r *AllSeriesRepository) FindAllByIds(c context.Context, ids []int) ([]models.AllSeries, error) {
	rows, err := r.db.Query(c, "select id, series,  title, description, release_year, director, rating, trailer_url from allseries where id = any($1)", ids)
	defer rows.Close()
	if err != nil {
		l := logger.GetLogger()
		l.Error(err.Error())
		return nil, err
	}

	allseries := make([]models.AllSeries, 0)

	for rows.Next() {
		var allserie models.AllSeries
		err = rows.Scan(&allserie.Id, &allserie.Series, &allserie.Title, &allserie.Description, &allserie.ReleaseYear, &allserie.Director, &allserie.Rating, &allserie.TrailerUrl)
		if err != nil {
			l := logger.GetLogger()
			l.Error(err.Error())
			return nil, err
		}

		allseries = append(allseries, allserie)
	}

	return allseries, nil
}

func (r *AllSeriesRepository) Create(c context.Context, serie models.AllSeries) (int, error) {
	var id int
	// tx, err := r.db.Begin(c)

	row := r.db.QueryRow(c, "insert into allseries (series,  title, description, release_year, director, rating, trailer_url) values($1, $2, $3, $4, $5, $6, $7) returning id", serie.Series)
	err := row.Scan(&id)
	if err != nil {
		l := logger.GetLogger()
		l.Error(err.Error())
		return 0, nil
	}
	return id, nil
}

func (r *AllSeriesRepository) FindById(c context.Context, id int) (models.AllSeries, error) {
	var allserie models.AllSeries
	row := r.db.QueryRow(c, "select id, series,  title, description, release_year, director, rating, trailer_url from allseries where id = $1", id)
	err := row.Scan(&allserie.Id,
		&allserie.Series,
		&allserie.Title,
		&allserie.Description,
		&allserie.ReleaseYear,
		&allserie.Director,
		&allserie.Rating,
		&allserie.TrailerUrl)
	if err != nil {
		l := logger.GetLogger()
		l.Error(err.Error())
		return models.AllSeries{}, err
	}
	return allserie, nil
}

func (r *AllSeriesRepository) FindAll(c context.Context) ([]models.AllSeries, error) {
	rows, err := r.db.Query(c, "select id, series, title, description, release_year, director, rating, trailer_url from allseries")
	defer rows.Close()
	if err != nil {
		l := logger.GetLogger()
		l.Error(err.Error())
		return nil, err
	}

	allseries := make([]models.AllSeries, 0)
	for rows.Next() {
		var allserie models.AllSeries
		err := rows.Scan(
			&allserie.Id,
			&allserie.Series,
			&allserie.Title,
			&allserie.Description,
			&allserie.ReleaseYear,
			&allserie.Director,
			&allserie.Rating,
			&allserie.TrailerUrl)
		if err != nil {
			l := logger.GetLogger()
			l.Error(err.Error())
			return nil, err
		}

		allseries = append(allseries, allserie)
	}

	return allseries, nil
}

func (r *AllSeriesRepository) Update(c context.Context, id int, allserie models.AllSeries) error {
	_, err := r.db.Exec(c, `update allseries set 
							series = $1 ,
							title = $2, 
							description = $3, 
							release_year = $4, 
							director = $5, 
							rating = $6, 
							trailer_url = $7
							where id = $8`,
		&allserie.Series,
		&allserie.Title,
		&allserie.Description,
		&allserie.ReleaseYear,
		&allserie.Director,
		&allserie.Rating,
		&allserie.TrailerUrl,
		id)
	if err != nil {
		l := logger.GetLogger()
		l.Error(err.Error())
		return err
	}

	return nil
}

//series, title, description, release_year, director, rating, trailer_url

func (r *AllSeriesRepository) Delete(c context.Context, id int) error {
	l := logger.GetLogger()

	var serieTitle string
	row := r.db.QueryRow(c, "select title from allseries where id = $1", id)
	err := row.Scan(&serieTitle)
	if err != nil {
		return err
	}

	l.Warn(fmt.Sprintf("Вы действительно хотите удалить %s серию?", serieTitle))

	_, err = r.db.Exec(c, "delete from allseries where id = $1", id)
	if err != nil {
		return err
	}

	l.Info(fmt.Sprintf("серия %s удалена", serieTitle))

	return nil
}
