package services

import (
	"errors"
	"penilaian_guru/dto"
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

func PatchVideoMetadata(db *gorm.DB, guruID uuid.UUID, req dto.PatchVideoRequest) (*models.VideoSubmission, error) {
	var video models.VideoSubmission
	err := db.Where("guru_id = ?", guruID).First(&video).Error
	if err != nil {
		return nil, errors.New("Video belum dikirim, tidak bisa diupdate")
	}

	if req.MataPelajaran != "" {
		video.MataPelajaran = req.MataPelajaran
	}
	if req.KelasSemester != "" {
		video.KelasSemester = req.KelasSemester
	}
	if req.HariTanggal != "" {
		video.HariTanggal = req.HariTanggal
	}
	if req.KompetensiDasar != "" {
		video.KompetensiDasar = req.KompetensiDasar
	}
	if req.Indikator != "" {
		video.Indikator = req.Indikator
	}

	video.UpdatedAt = time.Now()
	err = db.Save(&video).Error
	if err != nil {
		return nil, err
	}

	return &video, nil
}

func CreateVideoSubmission(
	db *gorm.DB,
	guruID uuid.UUID,
	link, mataPelajaran, kelasSemester, hariTanggal, kompetensiDasar, indikator string,
) (*models.VideoSubmission, error) {
	submission := models.VideoSubmission{
		ID:               uuid.New(),
		GuruID:           guruID,
		Link:             link,
		MataPelajaran:    mataPelajaran,
		KelasSemester:    kelasSemester,
		HariTanggal:      hariTanggal,
		KompetensiDasar:  kompetensiDasar,
		Indikator:        indikator,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(&submission).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func PatchVideoByID(db *gorm.DB, guruID, videoID uuid.UUID, req dto.PatchVideoRequest) (*models.VideoSubmission, error) {
	var video models.VideoSubmission
	err := db.Where("id = ? AND guru_id = ?", videoID, guruID).First(&video).Error
	if err != nil {
		return nil, errors.New("Video tidak ditemukan atau bukan milikmu")
	}

	if req.MataPelajaran != "" {
		video.MataPelajaran = req.MataPelajaran
	}
	if req.KelasSemester != "" {
		video.KelasSemester = req.KelasSemester
	}
	if req.HariTanggal != "" {
		video.HariTanggal = req.HariTanggal
	}
	if req.KompetensiDasar != "" {
		video.KompetensiDasar = req.KompetensiDasar
	}
	if req.Indikator != "" {
		video.Indikator = req.Indikator
	}

	video.UpdatedAt = time.Now()
	return &video, db.Save(&video).Error
}

type PenilaianGuruResponse struct {
	VideoID    uuid.UUID
	Link       string
	Label      string
	Catatan    string
	Saran      string
	KepsekNama string
}

func GetPenilaianGuru(db *gorm.DB, guruID uuid.UUID) ([]PenilaianGuruResponse, error) {
	var results []PenilaianGuruResponse

	err := db.Table("penilaians").
		Select(`
			penilaians.video_id,
			video_submissions.link,
			penilaians.label,
			penilaians.catatan,
			penilaians.saran,
			users.name as kepsek_nama
		`).
		Joins("JOIN video_submissions ON penilaians.video_id = video_submissions.id").
		Joins("JOIN users ON users.id = penilaians.kepsek_id").
		Where("video_submissions.guru_id = ?", guruID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

type VideoListItem struct {
	ID               uuid.UUID `json:"video_id"`
	Link             string    `json:"link"`
	MataPelajaran    string    `json:"mata_pelajaran"`
	KelasSemester    string    `json:"kelas_semester"`
	HariTanggal      string    `json:"hari_tanggal"`
	StatusPenilaian  string    `json:"status_penilaian"`
}

func GetAllVideosWithStatus(db *gorm.DB, guruID uuid.UUID) ([]VideoListItem, error) {
	var results []VideoListItem

	err := db.Model(&models.VideoSubmission{}).
		Select(`
			video_submissions.id,
			video_submissions.link,
			video_submissions.mata_pelajaran,
			video_submissions.kelas_semester,
			video_submissions.hari_tanggal,
			CASE
				WHEN p.id IS NOT NULL THEN 'Sudah dinilai'
				ELSE 'Menunggu untuk dinilai'
			END as status_penilaian
		`).
		Joins("LEFT JOIN penilaians p ON p.video_id = video_submissions.id").
		Where("video_submissions.guru_id = ?", guruID).
		Order("video_submissions.created_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetVideoByID(db *gorm.DB, guruID, videoID uuid.UUID) (*models.VideoSubmission, error) {
	var video models.VideoSubmission

	err := db.Where("id = ? AND guru_id = ?", videoID, guruID).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

type PenilaianDetail struct {
	Label       string
	Catatan     string
	Saran       string
	KepsekNama  string
	DinilaiPada time.Time
}

func GetPenilaianDetailByVideoID(db *gorm.DB, guruID, videoID uuid.UUID) (*PenilaianDetail, error) {
	var result PenilaianDetail

	err := db.Table("penilaians").
		Select(`
			penilaians.label,
			penilaians.catatan,
			penilaians.saran,
			kepsek.name as kepsek_nama,
			penilaians.created_at as dinilai_pada
		`).
		Joins("JOIN video_submissions ON penilaians.video_id = video_submissions.id").
		Joins("JOIN users kepsek ON penilaians.kepsek_id = kepsek.id").
		Where("penilaians.video_id = ? AND video_submissions.guru_id = ?", videoID, guruID).
		Scan(&result).Error

	if err != nil || result.KepsekNama == "" {
		return nil, errors.New("Video ini belum dinilai oleh kepala sekolah")
	}

	return &result, nil
}
