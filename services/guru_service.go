package services

import (
	"penilaian_guru/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpsertVideoSubmission(db *gorm.DB, guruID uuid.UUID, link string) (*models.VideoSubmission, error) {
	var submission models.VideoSubmission
	err := db.Where("guru_id = ?", guruID).First(&submission).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if submission.ID != uuid.Nil {

		submission.Link = link
		submission.UpdatedAt = time.Now()
		if err := db.Save(&submission).Error; err != nil {
			return nil, err
		}
		return &submission, nil
	}
	submission = models.VideoSubmission{
		ID:     uuid.New(),
		GuruID: guruID,
		Link:   link,
	}
	if err := db.Create(&submission).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}
