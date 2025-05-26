package services

import (
	"errors"
	"penilaian_guru/dto"
	"penilaian_guru/models"
	"penilaian_guru/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GuruSubmission struct {
	VideoID         uuid.UUID `json:"video_id"`
	GuruNama        string
	FotoProfil      string
	Link            string
	MataPelajaran   string
	KelasSemester   string
	HariTanggal     string
	UpdatedAt       string
	StatusPenilaian string `json:"status_penilaian"`
	Label           string `json:"label"`
}

func GetSubmissionsBySekolah(db *gorm.DB, sekolah string) ([]GuruSubmission, error) {
	var results []GuruSubmission

	err := db.Table("video_submissions").
		Select(`
		video_submissions.id as video_id,
		users.name as guru_nama,
		users.foto_profil,
		video_submissions.link,
		video_submissions.mata_pelajaran,
		video_submissions.kelas_semester,
		video_submissions.hari_tanggal,
		video_submissions.updated_at,
		CASE
			WHEN penilaians.id IS NOT NULL THEN 'Sudah dinilai'
			ELSE 'Menunggu untuk dinilai'
		END as status_penilaian`).
		Joins("JOIN users ON users.id = video_submissions.guru_id").
		Joins("LEFT JOIN penilaians ON penilaians.video_id = video_submissions.id").
		Where("users.sekolah = ?", sekolah).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetPenilaianByVideoID(db *gorm.DB, videoID uuid.UUID) (*models.Penilaian, error) {
	var penilaian models.Penilaian
	err := db.Where("video_id = ?", videoID).First(&penilaian).Error
	if err != nil {
		return nil, err
	}
	return &penilaian, nil
}

func CreatePenilaianByVideoID(db *gorm.DB, videoID, kepsekID uuid.UUID, input dto.PenilaianRequest) (*models.Penilaian, error) {
	var exists bool
	db.Model(&models.Penilaian{}).Select("count(*) > 0").
		Where("video_id = ?", videoID).Find(&exists)

	if exists {
		return nil, errors.New("Video ini sudah dinilai")
	}

	label := utils.HitungLabel(input)

	penilaian := models.Penilaian{
		ID:                 uuid.New(),
		VideoID:            videoID,
		KepsekID:           kepsekID,
		PersiapanMengajar:  input.PersiapanMengajar,
		MetodePembelajaran: input.MetodePembelajaran,
		PenguasaanMateri:   input.PenguasaanMateri,
		PengelolaanKelas:   input.PengelolaanKelas,
		Label:              label,
		Catatan:            input.Catatan,
		Saran:              input.Saran,
		CreatedAt:          time.Now(),
	}
	return &penilaian, db.Create(&penilaian).Error
}

func GetSubmissionByVideoID(db *gorm.DB, videoID uuid.UUID) (*GuruSubmission, error) {
	var result GuruSubmission

	err := db.Table("video_submissions").
		Select(`
			video_submissions.id as video_id,
			users.name as guru_nama,
			users.foto_profil,
			video_submissions.link,
			video_submissions.mata_pelajaran,
			video_submissions.kelas_semester,
			video_submissions.hari_tanggal,
			video_submissions.updated_at,
			COALESCE(p.label, 'Menunggu untuk dinilai') as label`).
		Joins("JOIN users ON users.id = video_submissions.guru_id").
		Joins("LEFT JOIN penilaians p ON p.video_id = video_submissions.id").
		Where("video_submissions.id = ?", videoID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}
