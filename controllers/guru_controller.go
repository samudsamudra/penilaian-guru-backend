package controllers

import (
	"net/http"
	"penilaian_guru/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoRequest struct {
	Link             string `json:"link" binding:"required,url"`
	MataPelajaran    string `json:"mata_pelajaran"`
	KelasSemester    string `json:"kelas_semester"`
	HariTanggal      string `json:"hari_tanggal"`
	KompetensiDasar  string `json:"kompetensi_dasar"`
	Indikator        string `json:"indikator"`
}

func SubmitVideoHandler(c *gin.Context, db *gorm.DB) {
	guruID := c.MustGet("userID").(uuid.UUID)

	var req VideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid."})
		return
	}

	submission, err := services.UpsertVideoSubmission(
		db,
		guruID,
		req.Link,
		req.MataPelajaran,
		req.KelasSemester,
		req.HariTanggal,
		req.KompetensiDasar,
		req.Indikator,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan video."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video berhasil dikirim!",
		"data": gin.H{
			"link":             submission.Link,
			"mata_pelajaran":   submission.MataPelajaran,
			"kelas_semester":   submission.KelasSemester,
			"hari_tanggal":     submission.HariTanggal,
			"kompetensi_dasar": submission.KompetensiDasar,
			"indikator":        submission.Indikator,
			"updated_at":       submission.UpdatedAt,
		},
	})
}

func GetMyVideoHandler(c *gin.Context, db *gorm.DB) {
	guruID := c.MustGet("userID").(uuid.UUID)

	video, err := services.GetVideoByGuruID(db, guruID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Belum ada video yang dikirim",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data video ditemukan",
		"data": gin.H{
			"link":              video.Link,
			"mata_pelajaran":    video.MataPelajaran,
			"kelas_semester":    video.KelasSemester,
			"hari_tanggal":      video.HariTanggal,
			"kompetensi_dasar":  video.KompetensiDasar,
			"indikator":         video.Indikator,
			"updated_at":        video.UpdatedAt,
		},
	})
}
