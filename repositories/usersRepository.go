package repositories

import (
	"context"
	"goozinshe/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(conn *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db: conn}
}

func (u *UsersRepository) FindById(c context.Context, id int) (models.User, error) {
	row := u.db.QueryRow(c, "select id, name, email, password_hash from users where id = $1", id)

	var user models.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)

	return user, err
}

func (u *UsersRepository) FindByEmail(c context.Context, email string) (models.User, error) {
	row := u.db.QueryRow(c, "select id, name, email, password_hash from users where email = $1", email)

	var user models.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)

	return user, err
}

func (r *UsersRepository) FindAll(c context.Context) ([]models.User, error) {
	rows, err := r.db.Query(c, "select id, name, email, password_hash from users order by id")
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) Create(c context.Context, user models.User) (int, error) {
	var id int
	err := r.db.QueryRow(c, "insert into users(name, email, password_hash) values($1, $2, $3) returning id", user.Name, user.Email, user.PasswordHash).Scan(&id)

	return id, err
}

func (u *UsersRepository) Update(c context.Context, id int, user models.User) error {
	_, err := u.db.Exec(c, "update users set name = $1, email = $2, password_hash = $3 where id = $4", user.Name, user.Email, user.PasswordHash, id)
	return err
}

func (u *UsersRepository) Delete(c context.Context, id int) error {
	_, err := u.db.Exec(c, "delete from users where id = $1", id)
	return err
}
