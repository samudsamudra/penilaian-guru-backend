package controllers

import (
	"context"
	"net/http"
	"os"
	"penilaian_guru/services"

	"penilaian_guru/utils"

	"github.com/gin-gonic/gin"
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
	user, created, err := services.FindOrCreateUser(db, userinfo.Email, userinfo.Name)
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

	msg := "Login berhasil"
	if created {
		msg = "Akun guru baru berhasil dibuat"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
		"token":   tokenString,
		"user":    user,
	})
}
