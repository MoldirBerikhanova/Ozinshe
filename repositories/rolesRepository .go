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

func (r *RolesRepository) FindById(c context.Context, id int) (models.Roles, error) {
	row := r.db.QueryRow(c, "select id, name, email, password_hash, phonenumber, birthday from roles where id = $1", id)

	var roles models.Roles
	err := row.Scan(&roles.Id, &roles.Name, &roles.Email, &roles.PasswordHash, &roles.PhoneNumber, &roles.Birthday)

	// PhoneNumber  int
	// Birthday     time.Tim

	return roles, err
}

func (r *RolesRepository) FindByEmail(c context.Context, email string) (models.Roles, error) {
	row := r.db.QueryRow(c, "select id, name, email, password_hash, phonenumber, birthday from roles where email = $1", email)

	var roles models.Roles
	err := row.Scan(&roles.Id, &roles.Name, &roles.Email, &roles.PasswordHash, &roles.PhoneNumber, &roles.Birthday)

	return roles, err
}

func (r *RolesRepository) FindAll(c context.Context) ([]models.Roles, error) {
	rows, err := r.db.Query(c, "select id, name, email, password_hash, phonenumber, birthday from roles order by id")
	if err != nil {
		return nil, err
	}

	roles := make([]models.Roles, 0)
	for rows.Next() {
		var role models.Roles
		err := rows.Scan(&role.Id, &role.Name, &role.Email, &role.PasswordHash, &role.PhoneNumber, &role.Birthday)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return roles, nil
}

func (r *RolesRepository) Create(c context.Context, role models.Roles) (int, error) {
	var id int
	err := r.db.QueryRow(c, "insert into roles(name, email, password_hash, phonenumber, birthday) values($1, $2, $3, $4, $5) returning id", role.Name, role.Email, role.PasswordHash, role.PhoneNumber, role.Birthday).Scan(&id)

	return id, err
}

func (r *RolesRepository) Update(c context.Context, id int, role models.Roles) error {
	_, err := r.db.Exec(c, "update roles set name = $1, email = $2, password_hash = $3, phonenumber = $4, birthday = $5  where id = $6", role.Name, role.Email, role.PasswordHash, role.PhoneNumber, role.Birthday, id)
	return err
}

func (r *RolesRepository) Delete(c context.Context, id int) error {
	_, err := r.db.Exec(c, "delete from roles where id = $1", id)
	return err
}
