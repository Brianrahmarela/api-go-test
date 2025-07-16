package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes: nama fungsi → dipakai untuk mendaftarkan semua route utama
// router *gin.Engine: object utama router Gin yang akan digunakan di main.go
// db *gorm.DB: koneksi ke database, diberikan agar bisa diteruskan ke controller
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// api :=: buat variabel bernama api
	// router.Group("/api"): semua route di bawah ini akan berada di prefix /api
	// Contoh hasil akhirnya: localhost:8080/api/register, localhost:8080/api/users
	api := router.Group("/api")
	{
		// Panggil fungsi dari file lain (auth_routes.go, user_routes.go)
		// Tujuannya: pisahkan logic auth & user supaya rapi
		SetupAuthRoutes(api, db)
		SetupUserRoutes(api, db)
	}
}
