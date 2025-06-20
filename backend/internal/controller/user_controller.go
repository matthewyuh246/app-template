package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"app-template/internal/entity"
	"app-template/internal/usecase"
)

// UserController ユーザーコントローラー
type UserController struct {
	userUseCase usecase.UserUseCase
}

// NewUserController ユーザーコントローラーの新しいインスタンスを作成
func NewUserController(userUseCase usecase.UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

// Register ユーザー登録ハンドラー
// @Summary ユーザー登録
// @Description 新しいユーザーを登録します
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.CreateUserRequest true "ユーザー登録リクエスト"
// @Success 201 {object} entity.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req entity.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	response, err := c.userUseCase.Register(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
			Code:  "REGISTRATION_ERROR",
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// Login ログインハンドラー
// @Summary ログイン
// @Description ユーザーログインを行います
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "ログインリクエスト"
// @Success 200 {object} entity.AuthResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req entity.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	response, err := c.userUseCase.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
			Code:  "LOGIN_ERROR",
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetUsers ユーザー一覧取得ハンドラー
// @Summary ユーザー一覧取得
// @Description ユーザー一覧を取得します
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(20)
// @Success 200 {object} entity.UsersResponse
// @Security BearerAuth
// @Router /users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	var params entity.PaginationParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	response, err := c.userUseCase.List(ctx.Request.Context(), &params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetUser ユーザー詳細取得ハンドラー
// @Summary ユーザー詳細取得
// @Description IDでユーザーの詳細を取得します
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Success 200 {object} entity.User
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	user, err := c.userUseCase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error: err.Error(),
			Code:  "USER_NOT_FOUND",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateUser ユーザー更新ハンドラー
// @Summary ユーザー更新
// @Description ユーザー情報を更新します
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Param request body entity.UpdateUserRequest true "ユーザー更新リクエスト"
// @Success 200 {object} entity.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	var req entity.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	user, err := c.userUseCase.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, ErrorResponse{
			Error: err.Error(),
			Code:  "UPDATE_ERROR",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser ユーザー削除ハンドラー
// @Summary ユーザー削除
// @Description ユーザーを削除します
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	err = c.userUseCase.Delete(ctx.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, ErrorResponse{
			Error: err.Error(),
			Code:  "DELETE_ERROR",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ErrorResponse エラーレスポンス
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
} 