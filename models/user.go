package models

import "time"
import "github.com/google/uuid"

type User struct {
    ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
    Role      string    `gorm:"not null"`
    Sekolah   string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Submissions []VideoSubmission `gorm:"foreignKey:GuruID"`
    FotoProfil string `gorm:"type:text"`
}
