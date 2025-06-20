package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"app-template/internal/entity"
)

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, id int64, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, params *entity.PaginationParams) ([]*entity.User, *entity.PaginationResponse, error)
}

// userRepository ユーザーリポジトリの実装
type userRepository struct {
	db *sql.DB
}

// NewUserRepository ユーザーリポジトリの新しいインスタンスを作成
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create 新しいユーザーを作成
func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
		INSERT INTO users (email, name, password, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`
	
	result, err := r.db.ExecContext(ctx, query, user.Email, user.Name, user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(ctx, id)
}

// GetByID IDでユーザーを取得
func (r *userRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, email, name, password, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// GetByEmail Emailでユーザーを取得
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, name, password, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// Update ユーザーを更新
func (r *userRepository) Update(ctx context.Context, id int64, user *entity.User) (*entity.User, error) {
	query := `
		UPDATE users 
		SET email = ?, name = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, user.Email, user.Name, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return r.GetByID(ctx, id)
}

// Delete ユーザーを削除
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List ユーザー一覧を取得
func (r *userRepository) List(ctx context.Context, params *entity.PaginationParams) ([]*entity.User, *entity.PaginationResponse, error) {
	// 総件数を取得
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count users: %w", err)
	}

	// ページネーション計算
	offset := (params.Page - 1) * params.Limit
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	// ユーザーリストを取得
	query := `
		SELECT id, email, name, password, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, params.Limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	pagination := &entity.PaginationResponse{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return users, pagination, nil
} 