package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"app-template/internal/entity"
	"app-template/internal/repository"
)

// UserUseCase ユーザーユースケースのインターフェース
type UserUseCase interface {
	Register(ctx context.Context, req *entity.CreateUserRequest) (*entity.AuthResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	Update(ctx context.Context, id int64, req *entity.UpdateUserRequest) (*entity.User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, params *entity.PaginationParams) (*entity.UsersResponse, error)
}

// userUseCase ユーザーユースケースの実装
type userUseCase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewUserUseCase ユーザーユースケースの新しいインスタンスを作成
func NewUserUseCase(userRepo repository.UserRepository, jwtSecret string) UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register 新しいユーザーを登録
func (u *userUseCase) Register(ctx context.Context, req *entity.CreateUserRequest) (*entity.AuthResponse, error) {
	// メールアドレスの重複チェック
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// ユーザーを作成
	user := &entity.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
	}

	createdUser, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// JWTトークンを生成
	token, err := u.generateJWT(createdUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &entity.AuthResponse{
		User:  createdUser,
		Token: token,
	}, nil
}

// Login ユーザーのログイン
func (u *userUseCase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error) {
	// ユーザーをメールアドレスで検索
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// パスワードを検証
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// JWTトークンを生成
	token, err := u.generateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &entity.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetByID IDでユーザーを取得
func (u *userUseCase) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// Update ユーザーを更新
func (u *userUseCase) Update(ctx context.Context, id int64, req *entity.UpdateUserRequest) (*entity.User, error) {
	// 既存ユーザーの確認
	existingUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	// メールアドレスの重複チェック（変更する場合）
	if req.Email != "" && req.Email != existingUser.Email {
		userWithEmail, err := u.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if userWithEmail != nil {
			return nil, fmt.Errorf("email already exists")
		}
	}

	// 更新データを準備
	updateUser := &entity.User{
		Email: existingUser.Email,
		Name:  existingUser.Name,
	}

	if req.Email != "" {
		updateUser.Email = req.Email
	}
	if req.Name != "" {
		updateUser.Name = req.Name
	}

	updatedUser, err := u.userRepo.Update(ctx, id, updateUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

// Delete ユーザーを削除
func (u *userUseCase) Delete(ctx context.Context, id int64) error {
	// ユーザーの存在確認
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	err = u.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List ユーザー一覧を取得
func (u *userUseCase) List(ctx context.Context, params *entity.PaginationParams) (*entity.UsersResponse, error) {
	// デフォルト値を設定
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	users, pagination, err := u.userRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return &entity.UsersResponse{
		Users:      users,
		Pagination: pagination,
	}, nil
}

// generateJWT JWTトークンを生成
func (u *userUseCase) generateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24時間有効
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
} 