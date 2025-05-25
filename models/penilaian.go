package models

import (
	"time"

	"github.com/google/uuid"
)

type Penilaian struct {
	ID                  uuid.UUID `gorm:"type:char(36);primaryKey"`
	VideoID             uuid.UUID `gorm:"type:char(36);not null;unique"`
	KepsekID            uuid.UUID `gorm:"type:char(36);not null"`
	PersiapanMengajar   int       `gorm:"not null"`
	MetodePembelajaran  int       `gorm:"not null"`
	PenguasaanMateri    int       `gorm:"not null"`
	PengelolaanKelas    int       `gorm:"not null"`
	SkorTotal           int       `gorm:"not null"`
	Label               string    `gorm:"not null"`
	Catatan             string    `gorm:"type:text"`
	Saran               string    `gorm:"type:text"`
	CreatedAt           time.Time

	Kepsek User `gorm:"foreignKey:KepsekID"`
}
