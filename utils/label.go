package utils

import "penilaian_guru/dto"

func HitungLabel(input dto.PenilaianRequest) string {
	total := input.PersiapanMengajar +
		input.MetodePembelajaran +
		input.PenguasaanMateri +
		input.PengelolaanKelas

	switch {
	case total >= 17:
		return "Sangat Baik"
	case total >= 13:
		return "Baik"
	case total >= 9:
		return "Cukup"
	default:
		return "Buruk"
	}
}
