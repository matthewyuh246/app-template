package entity

import (
	"time"
)

// User ユーザーエンティティ
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"` // JSONには含めない
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest ユーザー作成リクエスト
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest ユーザー更新リクエスト
type UpdateUserRequest struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
	Name  string `json:"name,omitempty" validate:"omitempty,min=1"`
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse 認証レスポンス
type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// PaginationParams ページネーションパラメータ
type PaginationParams struct {
	Page  int `form:"page" validate:"min=1"`
	Limit int `form:"limit" validate:"min=1,max=100"`
}

// PaginationResponse ページネーションレスポンス
type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// UsersResponse ユーザー一覧レスポンス
type UsersResponse struct {
	Users      []*User             `json:"users"`
	Pagination *PaginationResponse `json:"pagination"`
} 