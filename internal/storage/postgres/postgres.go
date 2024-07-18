package postgres

import (
	"cerberus/internal/domain/models"
	"cerberus/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

// NewStorage creates a new Storage instance.
func NewStorage(user, pass, name, host string) (*Storage, error) {
	const op = "internal/storage/postgres.NewStorage"

	db, err := sql.Open("postgres",
		fmt.Sprintf(
			"user=%s"+
				" password=%s"+
				" dbname=%s"+
				" host=%s"+
				" sslmode=disable",
			user,
			pass,
			name,
			host,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveUser saves a user to the database, handling both new user creation
// and existing user detection with improved error handling and logging.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (id int64, err error) {
	const op = "internal/storage/postgres.SaveUser"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// Check for existing user using a single prepared statement
	var existingID int64
	stmt, err := tx.PrepareContext(ctx, `
        SELECT 
            id
        FROM 
            users
        WHERE
            email = $1
    `)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRowContext(ctx, email).Scan(&existingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			stmt, err := tx.PrepareContext(ctx, `
                INSERT INTO users 
                    (email, password_hash) 
                VALUES
                    ($1, $2)
                RETURNING id
            `)
			if err != nil {
				return 0, fmt.Errorf("%s: %w", op, err)
			}

			err = stmt.QueryRowContext(ctx, email, passHash).Scan(&id)
			if err != nil {
				return 0, fmt.Errorf("%s: %w", op, err)
			}

			return id, nil
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Existing user detected, return error with custom type for clarity
	return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
}

// User retrieves a user from the database by email address.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "internal/storage/postgres.User"

	stmt, err := s.db.PrepareContext(ctx, `
		SELECT 
			id,
			email,
			password_hash
		FROM 
			users
		WHERE
			email = $1
	`)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	var user models.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// IsAdmin checks if a user is an admin.
func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "internal/storage/postgres.IsAdmin"

	stmt, err := s.db.PrepareContext(ctx, `
		SELECT 
			id
		FROM 
			admins
		WHERE
			user_id = $1
	`)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var isAdmin bool
	err = stmt.QueryRowContext(ctx, userId).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

// App retrieves an app from the database by ID.
func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "internal/storage/postgres.App"

	stmt, err := s.db.PrepareContext(ctx, `
		SELECT 
			id,
			name,
			secret
		FROM 
			apps
		WHERE
			id = $1
	`)
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	var app models.App
	err = stmt.QueryRowContext(ctx, appID).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
