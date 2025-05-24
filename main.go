package main

import (
    "penilaian_guru/config"
    "penilaian_guru/models"
    "penilaian_guru/routes"

    "github.com/gin-gonic/gin"
)

func main() {
    // Load environment (.env)
    config.LoadEnv()

    // Konek ke database
    db := config.ConnectDB()

    // Auto-migrate semua model ke database
    db.AutoMigrate(&models.User{}, &models.VideoSubmission{}, &models.Penilaian{})

    // Setup router
    r := gin.Default()
    routes.SetupRoutes(r, db)

    // Jalankan server di port 8080
    r.Run(":8080")
    
}
