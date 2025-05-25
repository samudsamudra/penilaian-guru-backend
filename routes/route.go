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
        guru.GET("/video", func(c *gin.Context) {
            controllers.GetMyVideoHandler(c, db)
        })
        guru.PATCH("/video", func(c *gin.Context) {
            controllers.PatchVideoHandler(c, db)
        })
        guru.PATCH("/video/:video_id", func(c *gin.Context) {
            controllers.PatchVideoByIDHandler(c, db)
        })
        guru.GET("/penilaian", func(c *gin.Context) {
            controllers.GetPenilaianHandler(c, db)
        })
    }

    r.GET("/me", middlewares.AuthMiddleware(), func(c *gin.Context) {
        controllers.GetMeHandler(c, db)
    })

    r.POST("/auth/kepsek/login", func(c *gin.Context) {
        controllers.KepsekLoginHandler(c, db)
    })

    r.GET("/auth/google/login", controllers.GoogleLogin)
    r.GET("/auth/google/callback", func(c *gin.Context) {
        controllers.GoogleCallback(c, db)
    })

    kepsek := r.Group("/kepsek")
    kepsek.Use(middlewares.AuthMiddleware())
    {
        kepsek.GET("/submissions", func(c *gin.Context) {
            controllers.GetSubmissionsByKepsekHandler(c, db)
        })
        kepsek.POST("/nilai", func(c *gin.Context) {
            controllers.PostPenilaianHandler(c, db)
        })
        kepsek.PATCH("/nilai/:video_id", func(c *gin.Context) {
            controllers.PatchPenilaianHandler(c, db)
        })
    }
}
