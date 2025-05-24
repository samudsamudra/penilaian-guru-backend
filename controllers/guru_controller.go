package controllers

import (
    "net/http"
    "penilaian_guru/services"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type SubmitVideoRequest struct {
    Link string `json:"link" binding:"required,url"`
}

func SubmitVideoHandler(c *gin.Context, db *gorm.DB) {
    guruID := c.MustGet("userID").(uuid.UUID)

    var req SubmitVideoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format link tidak valid."})
        return
    }

    submission, err := services.UpsertVideoSubmission(db, guruID, req.Link)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan video."})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Video berhasil dikirim!",
        "data": gin.H{
            "link":       submission.Link,
            "updated_at": submission.UpdatedAt,
        },
    })
}
