package models

import (
    "github.com/google/uuid"
    "time"
)

type VideoSubmission struct {
    ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
    GuruID    uuid.UUID `gorm:"type:char(36);not null"`
    Link      string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time

    MataPelajaran   string
    KelasSemester   string
    HariTanggal     string
    KompetensiDasar string
    Indikator       string


    Guru      User
    Penilaian Penilaian `gorm:"foreignKey:VideoID"`
}
