package routers

import (
	"blog-backend/controllers"
	"blog-backend/middleware"
	"blog-backend/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 使用全局日志+恢复中间件
	r.Use(middleware.Recovery(utils.Logger))

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		posts := api.Group("/posts")
		{
			posts.POST("/", controllers.CreatePost)
			posts.GET("/", controllers.GetPosts)
			posts.GET("/:id", controllers.GetPostByID)
			posts.PUT("/:id", controllers.UpdatePost)
			posts.DELETE("/:id", controllers.DeletePost)

			comments := posts.Group("/:id/comments")
			{
				comments.POST("", controllers.CreateComment)
				comments.GET("", controllers.GetCommentsByPost)
			}
		}
	}

	return r
}
