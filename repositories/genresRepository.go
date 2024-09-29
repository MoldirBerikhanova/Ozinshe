package repositories

import (
	"context"
	"errors"
	"goozinshe/models"
)

type GenresRepository struct {
	db map[int]models.Genre
}

func NewGenresRepository() *GenresRepository {
	return &GenresRepository{
		db: map[int]models.Genre{
			1: {
				Id:    1,
				Title: "Драма",
			},
			2: {
				Id:    2,
				Title: "Комедия",
			},
			3: {
				Id:    3,
				Title: "Ужасы",
			},
		},
	}
}

func (r *GenresRepository) FindById(c context.Context, id int) (models.Genre, error) {
	genre, ok := r.db[id]
	if !ok {
		return models.Genre{}, errors.New("genre not found")
	}

	return genre, nil
}

func (r *GenresRepository) FindAll(c context.Context) []models.Genre {
	genres := make([]models.Genre, 0, len(r.db))
	for _, v := range r.db {
		genres = append(genres, v)
	}

	return genres
}

func (r *GenresRepository) FindAllByIds(c context.Context, ids []int) []models.Genre {
	genres := make([]models.Genre, 0, len(r.db))
	for _, v := range r.db {
		for _, id := range ids {
			if v.Id == id {
				genres = append(genres, v)
			}
		}
	}

	return genres
}

func (r *GenresRepository) Create(c context.Context, genre models.Genre) int {
	id := len(r.db) + 1
	genre.Id = id

	r.db[id] = genre

	return genre.Id
}

func (r *GenresRepository) Update(c context.Context, id int, genre models.Genre) {
	original := r.db[id]
	original.Title = genre.Title

	r.db[id] = original
}

func (r *GenresRepository) Delete(c context.Context, id int) {
	delete(r.db, id)
}
