package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"app-template/internal/controller"
	"app-template/internal/repository"
	"app-template/internal/usecase"
	"app-template/pkg/database"
	"app-template/pkg/middleware"
)

// @title Web Application API
// @version 1.0
// @description Clean ArchitectureベースのWebアプリケーション用API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// データベース接続
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// リポジトリ層の初期化
	userRepo := repository.NewUserRepository(db)

	// ユースケース層の初期化
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key"
	}
	userUseCase := usecase.NewUserUseCase(userRepo, jwtSecret)

	// コントローラー層の初期化
	userController := controller.NewUserController(userUseCase)

	// Ginルーターの設定
	r := setupRouter(userController)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRouter ルーターの設定
func setupRouter(userController *controller.UserController) *gin.Engine {
	r := gin.Default()

	// ミドルウェアの設定
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": "2023-01-01T00:00:00Z",
		})
	})

	// API v1グループ
	v1 := r.Group("/api/v1")
	{
		// 認証関連（認証不要）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// ユーザー関連（認証必要）
		users := v1.Group("/users")
		users.Use(middleware.JWTAuth())
		{
			users.GET("", userController.GetUsers)
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}
	}

	return r
} 