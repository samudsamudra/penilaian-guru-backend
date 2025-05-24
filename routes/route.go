package routes

import (
    "penilaian_guru/controllers"

    "penilaian_guru/middlewares"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    guru := r.Group("/guru")
    guru.Use(middlewares.AuthMiddleware())
    {
        // Dummy: endpoint submit video
        guru.POST("/video", func(c *gin.Context) {
            controllers.SubmitVideoHandler(c, db)
        })
    }

    r.GET("/auth/google/login", controllers.GoogleLogin)
    r.GET("/auth/google/callback", func(c *gin.Context) {
        controllers.GoogleCallback(c, db)
    })
}
