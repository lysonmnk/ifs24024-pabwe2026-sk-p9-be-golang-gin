package server

import (
	"net/http"
	"time"

	"app/internal/config"
	"app/internal/features/ai"
	"app/internal/features/auth"
	"app/internal/features/post"
	"app/internal/middleware"
	"app/internal/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "app/docs"
)

func SetupRouter(
	cfg *config.Config,
	authH *auth.Handler,
	postH *post.Handler,
	aiH *ai.Handler,
) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORSMiddleware())

	limiter := middleware.NewRateLimiter(60, time.Minute)
	r.Use(limiter.Middleware())

	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "Endpoint not found", nil)
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/api.json", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/doc.json")
	})

	api := r.Group("/api")
	{
		// Public endpoints
		public := api.Group("/public")
		{
			public.GET("/categories", postH.GetCategories)
			public.GET("/posts", postH.GetPublicPosts)
			public.GET("/posts/:slug", postH.GetPublicPostBySlug)
		}

		// Auth endpoints
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authH.Register)
			authGroup.POST("/login", authH.Login)

			authProtected := authGroup.Group("")
			authProtected.Use(middleware.AuthMiddleware(cfg.AppKey))
			{
				authProtected.GET("/profile", authH.GetProfile)
				authProtected.PUT("/profile", authH.UpdateProfile)
				authProtected.POST("/logout", authH.Logout)
			}
		}

		// Author posts endpoints
		posts := api.Group("/posts")
		posts.Use(middleware.AuthMiddleware(cfg.AppKey))
		{
			posts.GET("", postH.GetAuthorPosts)
			posts.POST("", postH.CreatePost)
			posts.GET("/stats", postH.GetAuthorPostStats)
			posts.GET("/:id", postH.GetAuthorPostByID)
			posts.PUT("/:id", postH.UpdateAuthorPost)
			posts.PATCH("/:id/toggle-publish", postH.TogglePublish)
			posts.DELETE("/:id", postH.DeleteAuthorPost)
			posts.POST("/ai-generate", aiH.GeneratePost)
		}

		// Admin endpoints
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.AppKey))
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/users", postH.AdminGetAllUsers)
			admin.GET("/posts", postH.AdminGetAllPosts)
			admin.GET("/users/:id/posts", postH.AdminGetUserPosts)
		}
	}

	return r
}
