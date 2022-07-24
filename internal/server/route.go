package server

import (
	"restapi/internal/app/handler"
	"restapi/internal/app/repository"
	"restapi/internal/app/service"
	"restapi/internal/db/postgres"
	"restapi/internal/db/redis"
	"restapi/internal/logger"
	"restapi/internal/security/middleware"
	"restapi/internal/security/token"
	"restapi/internal/validation"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter(pg postgres.Client, rds redis.Client) *gin.Engine {
	router := gin.New()
	store := cookie.NewStore([]byte("secret"))
	router.Use(gin.Recovery(), logger.Logger(), middleware.CORSMiddleware(), sessions.Sessions("test", store))

	tk := token.NewToken()
	authRepo := repository.NewAuthRepo(rds)
	userRepo := repository.NewUserRepo(pg, rds)
	customRepo := repository.NewCustom(pg)

	authService := service.NewAuthService(userRepo, authRepo, tk)
	userService := service.NewUserService(userRepo, customRepo)

	authHandler := handler.NewAuthHandler(authService, tk)
	userHandler := handler.NewUserHandler(userService)

	api := router.Group("/api")
	api.POST("/login", authHandler.Login)
	api.POST("/register/:role", userHandler.Create)
	api.GET("/captcha", validation.CaptchaHandler)

	user := router.Group("/user", middleware.SetupAuthenticationMiddleware())
	user.GET("/:id", userHandler.Get)
	user.GET("/", userHandler.GetByToken)
	user.POST("/list", userHandler.List)
	user.PUT("/:id", middleware.Authorize(), userHandler.Update)
	user.PUT("/password/:id", middleware.Authorize(), userHandler.UpdatePassword)
	user.DELETE("/:id", middleware.Authorize(), userHandler.Delete)
	user.GET("/logout", authHandler.Logout)
	user.GET("/refresh", authHandler.Refresh)

	return router
}
