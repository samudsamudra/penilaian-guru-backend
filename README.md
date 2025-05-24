# Penilaian Guru Backend 🇮🇩

Sistem backend untuk platform penilaian kinerja guru berbasis video pembelajaran.  
Guru login menggunakan akun Google, lalu mengirimkan tautan video YouTube sebagai bahan supervisi.  
Kepala sekolah menilai melalui dashboard penilaian terintegrasi.

---

## 🔐 Fitur Utama

- **Google OAuth Login** untuk akun guru
- **JWT Authentication** untuk keamanan akses endpoint
- **Submit Video Pembelajaran** (link YouTube)
- **Middleware Auth** berbasis role (guru/kepsek/recon)
- **Auto-migrate** dengan GORM (MySQL)
- **Struktur clean-code** (modular controller, service, middleware)

---

## 🚀 Tech Stack

- Golang + Gin Gonic
- GORM + MySQL
- OAuth2 (Google)
- JWT (v5)
- PDFKit / Gofpdf (planned)
- GitHub CLI + VSCode

---

## 📂 Struktur Folder
```
penilaian_guru/
├── controllers/
├── services/
├── models/
├── routes/
├── middlewares/
├── utils/
├── config/
├── main.go
└── .env
```
---

## 🛠️ Setup

cp .env.example .env
go mod tidy
go run main.go
🧪 Endpoint Penting
Method	Endpoint	Deskripsi
GET	/ping	Cek server aktif
GET	/auth/google/login	Login guru via Google
GET	/auth/google/callback	Callback login
POST	/guru/video	Submit video (butuh token)

✨ Dev: Capellio Samudra
Made with Julian JJK
