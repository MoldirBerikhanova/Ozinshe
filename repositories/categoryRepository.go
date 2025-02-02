package repositories

import (
	"context"
	"goozinshe/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(conn *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: conn}
}

func (r *CategoryRepository) FindAllByIds(c context.Context, ids []int) ([]models.Category, error) {
	rows, err := r.db.Query(c, "select id, title from categories where id = any($1)", ids)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	categories := make([]models.Category, 0)

	for rows.Next() {
		var category models.Category
		err = rows.Scan(&category.Id, &category.Title)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) Create(c context.Context, category models.Category) (int, error) {
	var id int
	// tx, err := r.db.Begin(c)

	row := r.db.QueryRow(c, "insert into categories (title) values($1)", category.Title)
	err := row.Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (r *CategoryRepository) FindById(c context.Context, id int) (models.Category, error) {
	var category models.Category
	row := r.db.QueryRow(c, "select id, title from categories where id = $1", id)
	err := row.Scan(&category.Id, &category.Title)
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (r *CategoryRepository) FindAll(c context.Context) ([]models.Category, error) {
	rows, err := r.db.Query(c, "select id, title from categories")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Id, &category.Title)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) Update(c context.Context, id int, category models.Category) error {
	_, err := r.db.Exec(c, "update categories set title = $1 where id = $2", category.Title, category.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) Delete(c context.Context, id int) error {
	_, err := r.db.Exec(c, "delete from categories where id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
