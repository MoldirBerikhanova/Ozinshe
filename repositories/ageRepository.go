package repositories

import (
	"context"
	"goozinshe/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AgeRepository struct {
	db *pgxpool.Pool
}

func NewAgeRepository(conn *pgxpool.Pool) *AgeRepository {
	return &AgeRepository{db: conn}
}

func (r *AgeRepository) Create(c context.Context, age models.Age) (int, error) {
	var id int

	row := r.db.QueryRow(c, "insert into ages (title) values($1) returning id", age.Age)
	err := row.Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil

}

func (r *AgeRepository) FindAll(c context.Context) ([]models.Age, error) {
	rows, err := r.db.Query(c, "select id, age from ages")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	ages := make([]models.Age, 0)

	for rows.Next() {
		var age models.Age
		err := rows.Scan(&age.Id, &age.Age)
		if err != nil {
			return nil, err
		}
		ages = append(ages, age)
	}
	return ages, nil
}

func (r *AgeRepository) FindById(c context.Context, id int) (models.Age, error) {
	var age models.Age
	row := r.db.QueryRow(c, "select id, title from ages where id = $1", id)
	err := row.Scan(&age.Id, &age.Age)
	if err != nil {
		return models.Age{}, err
	}

	return age, nil
}

func (r *AgeRepository) FindAllByIds(c context.Context, ids []int) ([]models.Age, error) {
	rows, err := r.db.Query(c, "select id, age from ages where id = any($1)", ids)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	ages := make([]models.Age, 0)

	for rows.Next() {
		var age models.Age
		err := rows.Scan(&age.Id, &age.Age)
		if err != nil {
			return nil, err
		}
		ages = append(ages, age)
	}
	return ages, nil
}

func (r *AgeRepository) Update(c context.Context, id int, ages models.Age) error {
	_, err := r.db.Exec(c, "update ages set age = $1 where id = $2", ages.Age, ages.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *AgeRepository) Delete(c context.Context, id int) error {
	_, err := r.db.Exec(c, "delete from ages where id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
