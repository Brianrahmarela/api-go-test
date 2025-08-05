package controllers

import (
	"api-go-test/models"
	"net/http" // standar HTTP status code (seperti 404, 500)
	"strconv"  // konversi string ke angka

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileController struct {
	DB *gorm.DB
}

// constructor function
func NewProfileController(db *gorm.DB) *ProfileController {
	// Mmengembalikan pointer (&) ke struct yang di-isi properti DB-nya.
	return &ProfileController{DB: db}
}

// (pc *ProfileController) = method milik objek ProfileController.
// pc = receiver (seperti this di JS).
// *ProfileController = pointer ke struct ProfileController.
// c *gin.Context: objek context dari Gin, berisi permintaan & respons HTTP seperti body,
func (pc *ProfileController) CreateProfile(c *gin.Context) {
	// c.Param("id"): ambil :id dari URL (misal: /users/6 → "6").
	// strconv.Atoi(...): ubah string jadi integer ("6" → 6).
	// _ digunakan untuk mengabaikan error-nya (kalau gagal, akan diabaikan).
	userID, _ := strconv.Atoi(c.Param("id")) //id from user table

	var req models.CreateProfileRequest
	// VALIDASI BODY DARI CLIENT APAKAH JSON dgn c.ShouldBindJSON, jika tdk pakai ini, maka body tidak akan terbaca dan hanya terbaca sebagai model
	// c.ShouldBindJSON(&user): Baca body JSON dari request, masukkan ke struct user.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}
	//Membuat instance struct Profile berdasarkan data dari req.
	// UserID harus dikonversi ke uint karena model-nya pakai uint
	profile := models.Profile{
		UserID:    uint(userID), //convert to int
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
	}
	//create to table profile
	// VALIDASI SEBELUM CREATE
	//pc.DB.Create(&profile): simpan data ke database.
	//.Error: ambil error-nya, jika ada.
	if err := pc.DB.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create profile: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Profile created successfully",
		Data:    profile,
	})
}

func (pc *ProfileController) GetProfile(c *gin.Context) {
	userID := c.Param("id")
	var profile models.Profile

	pc.DB.Where("user_id = ?", userID).First(&profile)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    profile,
	})
}
