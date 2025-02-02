package repositories

import (
	"context"
	"goozinshe/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RolesRepository struct {
	db *pgxpool.Pool
}

func NewRolesRepository(conn *pgxpool.Pool) *RolesRepository {
	return &RolesRepository{db: conn}
}

func (r *RolesRepository) Create(c context.Context, roles models.Roles) (int, error) {
	var id int
	row := r.db.QueryRow(c, "insert into roles (names_of_hero, names_of_actors) values($1, $2)  returning id", roles.Names, roles.Actors)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RolesRepository) FindAll(c context.Context) ([]models.Roles, error) {
	rows, err := r.db.Query(c, "select  id, names_of_hero, names_of_actors from roles")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	roles := make([]models.Roles, 0)

	for rows.Next() {
		var role models.Roles
		err = rows.Scan(&role.Id, &role.Names, &role.Actors)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RolesRepository) FindAllByIds(c context.Context, ids []int) ([]models.Roles, error) {
	rows, err := r.db.Query(c, "select id, names_of_hero, names_of_actors from roles where id = any($1)", ids)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	roles := make([]models.Roles, 0)

	for rows.Next() {
		var role models.Roles
		err := rows.Scan(&role.Id, &role.Names, &role.Actors)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RolesRepository) FindById(c context.Context, id int) (models.Roles, error) {
	var role models.Roles
	row := r.db.QueryRow(c, "select id, names_of_hero, names_of_actors from roles where id = $1", id)
	err := row.Scan(&role.Id, &role.Names, &role.Actors)
	if err != nil {
		return models.Roles{}, err
	}
	return role, nil
}

func (r *RolesRepository) Update(c context.Context, id int, updatedroles models.Roles) error {
	_, err := r.db.Exec(c, "update roles set names_of_hero = $1,  names_of_actors = $2  where id = $3", updatedroles.Id, updatedroles.Names, updatedroles.Actors)
	if err != nil {
		return err
	}

	return nil
}

func (r *RolesRepository) Delete(c context.Context, id int) error {
	_, err := r.db.Exec(c, "delete from roles where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
