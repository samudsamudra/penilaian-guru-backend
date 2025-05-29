package controllers

import (
	"context"
	"net/http"
	"os"
	"penilaian_guru/services"

	"penilaian_guru/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin(c *gin.Context) {
	conf := getGoogleOAuthConfig()
	url := conf.AuthCodeURL("randomstate")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context, db *gorm.DB) {
	conf := getGoogleOAuthConfig()

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}
	oauth2Service, err := oauth2api.NewService(context.Background(), option.WithTokenSource(conf.TokenSource(context.Background(), token)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create oauth2 service"})
		return
	}

	userinfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Buat atau ambil user dari database
	user, _, err := services.FindOrCreateUser(db, userinfo.Email, userinfo.Name, userinfo.Picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal login / register user"})
		return
	}

	// Buat JWT token
	tokenString, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	// Kirim token via Cookie
	c.SetCookie(
		"token",         // nama cookie
		tokenString,     // isi token
		3600,            // durasi (1 jam)
		"/",             // path
		"localhost",     // domain (sesuaikan kalau nanti pakai domain)
		false,           // secure: false di localhost, true kalau HTTPS
		true,            // httpOnly: true supaya FE gak bisa akses langsung dari JS
	)

    c.SetCookie(
        "guruID",       // nama cookie
        user.ID.String(), // isi user ID
        3600,            // durasi (1 jam)
        "/",             // path
        "localhost",     // domain (sesuaikan kalau nanti pakai domain)
        false,           // secure: false di localhost, true kalau HTTPS
        true,            // httpOnly: true supaya FE gak bisa akses langsung dari JS
    )

	// Redirect ke dashboard FE
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")
}

func GetMeHandler(c *gin.Context, db *gorm.DB) {
	userID := c.MustGet("userID").(uuid.UUID)

	user, err := services.GetUserByID(db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data user ditemukan",
		"user": gin.H{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"role":        user.Role,
			"sekolah":     user.Sekolah,
			"foto_profil": user.FotoProfil,
		},
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func KepsekLoginHandler(c *gin.Context, db *gorm.DB) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format login tidak valid"})
		return
	}

	key := req.Email // or use c.ClientIP()

	if utils.IsBlocked(key) {
		c.JSON(429, gin.H{"error": "jangan bruteforce pliss, tunggu beberapa saat lagi"})
		return
	}

	user, err := services.GetUserByEmail(db, req.Email)
	if err != nil || user.Role != "kepsek" {
		utils.RegisterFail(key)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akun tidak ditemukan atau bukan kepsek"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.RegisterFail(key)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	utils.ResetAttempts(key)

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

    // Set cookie untuk token
    c.SetCookie(
        "token",
        token,
        3600,
        "/",
        "localhost",
        false,
        true,
    )

    // Set cookie untuk user ID
    c.SetCookie(
        "userID",
        user.ID.String(),
        3600,
        "/",
        "localhost",
        false,
        true,
    )

	c.JSON(http.StatusOK, gin.H{
		"message": "Login kepsek berhasil",
		"token":   token,
		"user": gin.H{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"role":    user.Role,
			"sekolah": user.Sekolah,
		},
	})
}
