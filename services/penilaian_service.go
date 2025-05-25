package services

import (
	"errors"
	"penilaian_guru/dto"
	"penilaian_guru/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func LabelDariSkor(skor int) string {
	switch {
	case skor >= 18:
		return "Sangat Baik"
	case skor >= 15:
		return "Baik"
	case skor >= 10:
		return "Cukup"
	default:
		return "Buruk"
	}
}

func UpsertPenilaian(
	db *gorm.DB,
	videoID uuid.UUID,
	kepsekID uuid.UUID,
	persiapan, metode, materi, kelas int,
	catatan, saran string,
) (*models.Penilaian, error) {
	total := persiapan + metode + materi + kelas
	label := LabelDariSkor(total)

	var existing models.Penilaian
	err := db.Where("video_id = ?", videoID).First(&existing).Error
	if err == nil {
		// update
		existing.PersiapanMengajar = persiapan
		existing.MetodePembelajaran = metode
		existing.PenguasaanMateri = materi
		existing.PengelolaanKelas = kelas
		existing.Catatan = catatan
		existing.Saran = saran
		existing.SkorTotal = total
		existing.Label = label
		return &existing, db.Save(&existing).Error
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// insert baru
	penilaian := models.Penilaian{
		ID:                 uuid.New(),
		VideoID:            videoID,
		KepsekID:           kepsekID,
		PersiapanMengajar:  persiapan,
		MetodePembelajaran: metode,
		PenguasaanMateri:   materi,
		PengelolaanKelas:   kelas,
		Catatan:            catatan,
		Saran:              saran,
		SkorTotal:          total,
		Label:              label,
		CreatedAt:          time.Now(),
	}
	return &penilaian, db.Create(&penilaian).Error
}

func PatchPenilaian(
	db *gorm.DB,
	videoID uuid.UUID,
	kepsekID uuid.UUID,
	req dto.PatchPenilaianRequest,
) (*models.Penilaian, error) {
	var p models.Penilaian

	err := db.Where("video_id = ? AND kepsek_id = ?", videoID, kepsekID).First(&p).Error
	if err != nil {
		return nil, errors.New("penilaian tidak ditemukan")
	}

	// Update field jika dikirim
	if req.PersiapanMengajar != nil {
		p.PersiapanMengajar = *req.PersiapanMengajar
	}
	if req.MetodePembelajaran != nil {
		p.MetodePembelajaran = *req.MetodePembelajaran
	}
	if req.PenguasaanMateri != nil {
		p.PenguasaanMateri = *req.PenguasaanMateri
	}
	if req.PengelolaanKelas != nil {
		p.PengelolaanKelas = *req.PengelolaanKelas
	}
	if req.Catatan != "" {
		p.Catatan = req.Catatan
	}
	if req.Saran != "" {
		p.Saran = req.Saran
	}

	// Hitung ulang total & label
	total := p.PersiapanMengajar + p.MetodePembelajaran + p.PenguasaanMateri + p.PengelolaanKelas
	p.SkorTotal = total
	p.Label = LabelDariSkor(total)

	return &p, db.Save(&p).Error
}
