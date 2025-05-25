package services

import (
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
