package services

import (
	"errors"
	"penilaian_guru/dto"
	"penilaian_guru/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpsertVideoSubmission(
	db *gorm.DB,
	guruID uuid.UUID,
	link, mataPelajaran, kelasSemester, hariTanggal, kompetensiDasar, indikator string,
) (*models.VideoSubmission, error) {
	var submission models.VideoSubmission
	err := db.Where("guru_id = ?", guruID).First(&submission).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	now := time.Now()

	if err == nil {
		// update existing
		submission.Link = link
		submission.MataPelajaran = mataPelajaran
		submission.KelasSemester = kelasSemester
		submission.HariTanggal = hariTanggal
		submission.KompetensiDasar = kompetensiDasar
		submission.Indikator = indikator
		submission.UpdatedAt = now

		if err := db.Save(&submission).Error; err != nil {
			return nil, err
		}
		return &submission, nil
	}

	// new submission
	submission = models.VideoSubmission{
		ID:               uuid.New(),
		GuruID:           guruID,
		Link:             link,
		MataPelajaran:    mataPelajaran,
		KelasSemester:    kelasSemester,
		HariTanggal:      hariTanggal,
		KompetensiDasar:  kompetensiDasar,
		Indikator:        indikator,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := db.Create(&submission).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}
func GetVideoByGuruID(db *gorm.DB, guruID uuid.UUID) (*models.VideoSubmission, error) {
	var submission models.VideoSubmission
	err := db.Where("guru_id = ?", guruID).First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

func PatchVideoMetadata(db *gorm.DB, guruID uuid.UUID, req dto.PatchVideoRequest) (*models.VideoSubmission, error) {
	var video models.VideoSubmission
	err := db.Where("guru_id = ?", guruID).First(&video).Error
	if err != nil {
		return nil, errors.New("Video belum dikirim, tidak bisa diupdate")
	}

	if req.MataPelajaran != "" {
		video.MataPelajaran = req.MataPelajaran
	}
	if req.KelasSemester != "" {
		video.KelasSemester = req.KelasSemester
	}
	if req.HariTanggal != "" {
		video.HariTanggal = req.HariTanggal
	}
	if req.KompetensiDasar != "" {
		video.KompetensiDasar = req.KompetensiDasar
	}
	if req.Indikator != "" {
		video.Indikator = req.Indikator
	}

	video.UpdatedAt = time.Now()
	err = db.Save(&video).Error
	if err != nil {
		return nil, err
	}

	return &video, nil
}

func CreateVideoSubmission(
	db *gorm.DB,
	guruID uuid.UUID,
	link, mataPelajaran, kelasSemester, hariTanggal, kompetensiDasar, indikator string,
) (*models.VideoSubmission, error) {
	submission := models.VideoSubmission{
		ID:               uuid.New(),
		GuruID:           guruID,
		Link:             link,
		MataPelajaran:    mataPelajaran,
		KelasSemester:    kelasSemester,
		HariTanggal:      hariTanggal,
		KompetensiDasar:  kompetensiDasar,
		Indikator:        indikator,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(&submission).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func PatchVideoByID(db *gorm.DB, guruID, videoID uuid.UUID, req dto.PatchVideoRequest) (*models.VideoSubmission, error) {
	var video models.VideoSubmission
	err := db.Where("id = ? AND guru_id = ?", videoID, guruID).First(&video).Error
	if err != nil {
		return nil, errors.New("Video tidak ditemukan atau bukan milikmu")
	}

	if req.MataPelajaran != "" {
		video.MataPelajaran = req.MataPelajaran
	}
	if req.KelasSemester != "" {
		video.KelasSemester = req.KelasSemester
	}
	if req.HariTanggal != "" {
		video.HariTanggal = req.HariTanggal
	}
	if req.KompetensiDasar != "" {
		video.KompetensiDasar = req.KompetensiDasar
	}
	if req.Indikator != "" {
		video.Indikator = req.Indikator
	}

	video.UpdatedAt = time.Now()
	return &video, db.Save(&video).Error
}

type PenilaianGuruResponse struct {
	VideoID    uuid.UUID
	Link       string
	Label      string
	Catatan    string
	Saran      string
	KepsekNama string
}

func GetPenilaianGuru(db *gorm.DB, guruID uuid.UUID) ([]PenilaianGuruResponse, error) {
	var results []PenilaianGuruResponse

	err := db.Table("penilaians").
		Select(`
			penilaians.video_id,
			video_submissions.link,
			penilaians.label,
			penilaians.catatan,
			penilaians.saran,
			users.name as kepsek_nama
		`).
		Joins("JOIN video_submissions ON penilaians.video_id = video_submissions.id").
		Joins("JOIN users ON users.id = penilaians.kepsek_id").
		Where("video_submissions.guru_id = ?", guruID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
