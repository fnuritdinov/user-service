package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/fnuritdinov/user-service/internal/models"
	errs "github.com/fnuritdinov/user-service/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return Repository{db: db}
}

func (r *Repository) Register(ctx context.Context, request models.User) (int, error) {
	const query = `
			INSERT INTO mv_users (name, email, password, phone, age, role)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`

	var id int
	err := r.db.QueryRow(
		ctx, query, request.Name, request.Email, request.Password, request.Phone, request.Age, request.Role).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `
			SELECT EXISTS (
				SELECT 1 
				FROM users
				WHERE email = $1
			);`

	var exist bool

	err := r.db.QueryRow(ctx, query, email).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var myUser models.User

	const query = `SELECT id, password, role FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&myUser.ID,
		&myUser.Email,
		&myUser.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errs.ErrUserNotFound
		}
		return models.User{}, err
	}

	return myUser, nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, req models.HashTokenReq) error {
	return nil
}

func (r *Repository) GetRefreshTokenByHash(ctx context.Context, hash string) (models.HashTokenReq, error) {
	var token models.HashTokenReq
	const query = `SELECT id, user_id, expired_at FROM hash_tokens WHERE hash = $1`

	err := r.db.QueryRow(ctx, query, hash).
		Scan(&token.ID,
			&token.UserID,
			&token.ExpiredAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.HashTokenReq{}, errs.ErrNotFound
		}

		return models.HashTokenReq{}, fmt.Errorf("error from r.db.QueryRow %w", err)
	}

	return token, nil
}

func (r *Repository) DeleteRefreshTokenByID(ctx context.Context, tokenID int) error {
	const query = `DELETE FROM hash_tokens WHERE id = $1`

	result, err := r.db.Exec(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("error from delete token %w", err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
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

	const query = `UPDATE mv_users SET name = $2, phone = $3 WHERE id = $1`

	result, err := r.db.Exec(ctx, query, request.ID, request.Name, request.Phone)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *Repository) GetList(ctx context.Context) ([]models.User, error) {
	const query = `SELECT * FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []models.User{}, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.Age,
			&user.Name,
			&user.Role,
			&user.Phone,
			&user.Password)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
