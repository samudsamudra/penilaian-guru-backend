package services

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GuruSubmission struct {
	VideoID         uuid.UUID       `json:"video_id"`
	GuruNama        string
	FotoProfil      string
	Link            string
	MataPelajaran   string
	KelasSemester   string
	HariTanggal     string
	UpdatedAt       string
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
		video_submissions.updated_at`).
	Joins("JOIN users ON users.id = video_submissions.guru_id").
	Where("users.sekolah = ?", sekolah).
	Scan(&results).Error



	if err != nil {
		return nil, err
	}
	return results, nil
}
