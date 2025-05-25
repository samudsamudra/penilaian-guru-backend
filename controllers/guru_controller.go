package controllers

import (
	"net/http"
	"penilaian_guru/dto"
	"penilaian_guru/services"

	// Added import
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

	submission, err := services.CreateVideoSubmission(
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

func PatchVideoHandler(c *gin.Context, db *gorm.DB) {
	guruID := c.MustGet("userID").(uuid.UUID)

	var req dto.PatchVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	video, err := services.PatchVideoMetadata(db, guruID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Metadata video berhasil diperbarui",
		"data":    video,
	})
}

func PatchVideoByIDHandler(c *gin.Context, db *gorm.DB) {
	guruID := c.MustGet("userID").(uuid.UUID)
	videoIDParam := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID video tidak valid"})
		return
	}

	var req dto.PatchVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	video, err := services.PatchVideoByID(db, guruID, videoID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Metadata video berhasil diperbarui",
		"data":    video,
	})
}

func GetPenilaianHandler(c *gin.Context, db *gorm.DB) {
	guruID := c.MustGet("userID").(uuid.UUID)

	results, err := services.GetPenilaianGuru(db, guruID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil penilaian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar penilaian ditemukan",
		"data":    results,
	})
}

func GetAllVideoHandler(c *gin.Context, db *gorm.DB) {
	guruID := c.MustGet("userID").(uuid.UUID)

	videos, err := services.GetAllVideosWithStatus(db, guruID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar video ditemukan",
		"data":    videos,
	})
}

func GetVideoDetailHandler(c *gin.Context, db *gorm.DB) {
	guruID := c.MustGet("userID").(uuid.UUID)
	videoIDStr := c.Param("id")

	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	video, err := services.GetVideoByID(db, guruID, videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video tidak ditemukan atau bukan milikmu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail video ditemukan",
		"data": gin.H{
			"video_id":         video.ID,
			"link":             video.Link,
			"mata_pelajaran":   video.MataPelajaran,
			"kelas_semester":   video.KelasSemester,
			"hari_tanggal":     video.HariTanggal,
			"kompetensi_dasar": video.KompetensiDasar,
			"indikator":        video.Indikator,
			"updated_at":       video.UpdatedAt,
		},
	})
}

func GetPenilaianDetailHandler(c *gin.Context, db *gorm.DB) {
	guruID := c.MustGet("userID").(uuid.UUID)

	idStr := c.Param("id")
	videoID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result, err := services.GetPenilaianDetailByVideoID(db, guruID, videoID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Video ini belum dinilai",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penilaian ditemukan",
		"data": gin.H{
			"label":        result.Label,
			"catatan":      result.Catatan,
			"saran":        result.Saran,
			"kepsek_nama":  result.KepsekNama,
			"dinilai_pada": result.DinilaiPada,
		},
	})
}
