package dto

type PatchPenilaianRequest struct {
	PersiapanMengajar  *int   `json:"persiapan_mengajar"`
	MetodePembelajaran *int   `json:"metode_pembelajaran"`
	PenguasaanMateri   *int   `json:"penguasaan_materi"`
	PengelolaanKelas   *int   `json:"pengelolaan_kelas"`
	Catatan            string `json:"catatan"`
	Saran              string `json:"saran"`
}

type PenilaianRequest struct {
	PersiapanMengajar  int    `json:"persiapan_mengajar"`
	MetodePembelajaran int    `json:"metode_pembelajaran"`
	PenguasaanMateri   int    `json:"penguasaan_materi"`
	PengelolaanKelas   int    `json:"pengelolaan_kelas"`
	Catatan            string `json:"catatan"`
	Saran              string `json:"saran"`
}
