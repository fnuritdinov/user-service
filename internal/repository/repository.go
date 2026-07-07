package repository

import (
	"context"
	"errors"
	"user-service/internal/models"
	errs "user-service/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, request models.User) (int, error) {
	const query = `
			INSERT INTO mv_users (name, email, password, phone, age)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`

	var id int
	err := r.db.QueryRow(ctx, query, request.Name, request.Email, request.Password, request.Phone, request.Age).Scan(&id)
	if err != nil {
		return 0, errors.New("error from r.db.Exec")
	}

	return id, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (models.User, error) {
	const query = `SELECT id, name, age, email, password, phone FROM mv_users WHERE id = $1`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.Email,
		&user.Password,
		&user.Phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errs.ErrUserNotFound
		}

		return models.User{}, err
	}
	return user, err

}

func (r *Repository) Update(ctx context.Context, request models.UpdateUser) error {

	const query = `UPDATE mv_users SET name = $2, phone = $3, WHERE id = $1`

	result, err := r.db.Exec(ctx, query, request.ID, request.Name, request.Phone)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}
