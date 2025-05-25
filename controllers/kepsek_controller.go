package controllers

import (
	"net/http"
	"penilaian_guru/dto"
	"penilaian_guru/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetSubmissionsByKepsekHandler(c *gin.Context, db *gorm.DB) {
	userID := c.MustGet("userID").(uuid.UUID)

	user, err := services.GetUserByID(db, userID)
	if err != nil || user.Role != "kepsek" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak"})
		return
	}

	results, err := services.GetSubmissionsBySekolah(db, user.Sekolah)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar video guru berhasil diambil",
		"data":    results,
	})
}

type PenilaianRequest struct {
	VideoID      uuid.UUID         `json:"video_id"`
	SkorPerAspek map[string]int    `json:"skor_per_aspek"`
}

type PenilaianRequestBaru struct {
	VideoID            uuid.UUID `json:"video_id"`
	PersiapanMengajar  int       `json:"persiapan_mengajar"`
	MetodePembelajaran int       `json:"metode_pembelajaran"`
	PenguasaanMateri   int       `json:"penguasaan_materi"`
	PengelolaanKelas   int       `json:"pengelolaan_kelas"`
	Catatan            string    `json:"catatan"`
	Saran              string    `json:"saran"`
}

func PostPenilaianHandler(c *gin.Context, db *gorm.DB) {
	kepsekID := c.MustGet("userID").(uuid.UUID)

	var req PenilaianRequestBaru
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format penilaian tidak valid"})
		return
	}

	penilaian, err := services.UpsertPenilaian(
		db,
		req.VideoID,
		kepsekID,
		req.PersiapanMengajar,
		req.MetodePembelajaran,
		req.PenguasaanMateri,
		req.PengelolaanKelas,
		req.Catatan,
		req.Saran,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan penilaian"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penilaian berhasil disimpan",
		"data": gin.H{
			"total_skor": penilaian.SkorTotal,
			"label":      penilaian.Label,
		},
	})
}

type PatchPenilaianRequest struct {
	PersiapanMengajar  *int   `json:"persiapan_mengajar"`
	MetodePembelajaran *int   `json:"metode_pembelajaran"`
	PenguasaanMateri   *int   `json:"penguasaan_materi"`
	PengelolaanKelas   *int   `json:"pengelolaan_kelas"`
	Catatan            string `json:"catatan"`
	Saran              string `json:"saran"`
}

func PatchPenilaianHandler(c *gin.Context, db *gorm.DB) {
	kepsekID := c.MustGet("userID").(uuid.UUID)
	videoIDParam := c.Param("video_id")
	videoID, err := uuid.Parse(videoIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID video tidak valid"})
		return
	}

	var req PatchPenilaianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	penilaian, err := services.PatchPenilaian(db, videoID, kepsekID, dto.PatchPenilaianRequest(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penilaian berhasil diperbarui",
		"data": gin.H{
			"total_skor": penilaian.SkorTotal,
			"label":      penilaian.Label,
		},
	})
}
