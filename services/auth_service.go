package services

import (
	"penilaian_guru/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindOrCreateUser(db *gorm.DB, email, name, foto string) (*models.User, bool, error) {
    var user models.User
    err := db.Where("email = ?", email).First(&user).Error
    if err == nil {
        return &user, false, nil // udah ada
    }
    if err != gorm.ErrRecordNotFound {
        return nil, false, err // error lain
    }

    user = models.User{
        ID:      uuid.New(),
        Email:   email,
        Name:    name,
        Role:    "guru",
        Sekolah: "SMK Telkom Malang",
        FotoProfil: foto,
    }

    if err := db.Create(&user).Error; err != nil {
        return nil, false, err
    }
    return &user, true, nil
}

func GetUserByID(db *gorm.DB, userID uuid.UUID) (*models.User, error) {
    var user models.User
    err := db.First(&user, "id = ?", userID).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
    var user models.User
    err := db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
