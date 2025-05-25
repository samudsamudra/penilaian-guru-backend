package dto

type PatchVideoRequest struct {
	MataPelajaran    string `json:"mata_pelajaran"`
	KelasSemester    string `json:"kelas_semester"`
	HariTanggal      string `json:"hari_tanggal"`
	KompetensiDasar  string `json:"kompetensi_dasar"`
	Indikator        string `json:"indikator"`
}
