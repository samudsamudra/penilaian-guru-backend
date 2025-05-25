package models

import (
	"time"

	"github.com/google/uuid"
)

type VideoSubmission struct {
    ID               uuid.UUID `gorm:"type:char(36);primaryKey"`
    GuruID           uuid.UUID `gorm:"type:char(36);not null"`
    Link             string    `gorm:"not null"`
    MataPelajaran    string
    KelasSemester    string
    HariTanggal      string
    KompetensiDasar  string
    Indikator        string
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
