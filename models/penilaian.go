package models

import (
    "github.com/google/uuid"
    "time"
)

type Penilaian struct {
    ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
    VideoID      uuid.UUID `gorm:"type:char(36);not null;unique"`
    KepsekID     uuid.UUID `gorm:"type:char(36);not null"`
    SkorTotal    int       `gorm:"not null"`
    SkorPerAspek string    `gorm:"type:text"`
    Label        string    `gorm:"not null"`
    CreatedAt    time.Time

    Kepsek User `gorm:"foreignKey:KepsekID"`
}
