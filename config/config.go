package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("❗Warning: .env file not found. Using default env.")
    }
}

func ConnectDB() *gorm.DB {
    dsn := os.Getenv("DB_DSN") // contoh: "user:password@tcp(127.0.0.1:3306)/dbname?parseTime=true"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("❌ Gagal konek DB:", err)
    }
    fmt.Println("✅ Connected to DB!")
    DB = db
    return db
}
