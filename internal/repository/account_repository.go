package repository

import (
	"context"
	"database/sql"
	"errors"

	"bank-app/internal/model"
)

type AccountRepo struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) Create(ctx context.Context, account *model.Account) error {
	query := `
		INSERT INTO accounts (user_id, number, balance, currency)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		account.UserID,
		account.Number,
		account.Balance,
		account.Currency,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *AccountRepo) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	account := &model.Account{}
	query := `
		SELECT id, user_id, number, balance, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Number,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("account not found")
	}

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepo) GetByUserID(ctx context.Context, userID int64) ([]*model.Account, error) {
	query := `
		SELECT id, user_id, number, balance, currency, created_at, updated_at
		FROM accounts
		WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*model.Account
	for rows.Next() {
		account := &model.Account{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Number,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepo) Update(ctx context.Context, account *model.Account) error {
	query := `
		UPDATE accounts
		SET balance = $1, currency = $2
		WHERE id = $3
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		account.Balance,
		account.Currency,
		account.ID,
	).Scan(&account.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}
