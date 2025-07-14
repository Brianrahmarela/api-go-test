package controllers

import (
	"api-go-test/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Menyimpan logic yang akan merespon HTTP request
type UserController struct {
	//UserService: pointer ke UserService dari package services
	//Tujuannya: bisa panggil fungsi di service utk ambil user dari DB
	UserService *services.UserService
}

//gambaran bentuk blueprint struct UserController
// UserController {
//     UserService --> (pointer ke struct UserService)
// }

// Buat controller + isi UserService-nya
// Constructor hanya fungsi biasa yang mengembalikan struct
func NewUserController(db *gorm.DB) *UserController {
	// Bungkus ke dalam struct UserController, trs return alamatnya (&UserController) ke function
	return &UserController{
		//Buat object UserService lewat services.NewUserService(db)
		UserService: services.NewUserService(db),
	}
}

//gambaran bentuk blueprint struct UserController (objek yang dibuat):
// &UserController {
//     UserService: &UserService {
//         DB: &gorm.DB  → koneksi database aktif
//     }
// }

// uc *UserController → method ini milik UserController bernama GetUsers
// GetUsers() -> Mengambil semua user dari service lalu kirim JSON
// GetUsers bukan property atau field yang tersimpan di dalam UserController
// Tapi dia adalah fungsi (method) yang “nempel” ke tipe *UserController
// method bisa "terikat" ke suatu struct type, Dihubungkan lewat receiver (contohnya func (uc *UserController) ...)
func (uc *UserController) GetUsers(c *gin.Context) {
	//Panggil GetAllUsers() dari UserService
	users, err := uc.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	//Jika sukses, kirim response 200 (OK) dengan data user-nya
	c.JSON(http.StatusOK, gin.H{"data": users})
}

//gambaran bentuk method yg nempel ke struct UserController
// func (uc *UserController) GetUsers(...)

func (uc *UserController) GetUserByID(c *gin.Context) {
	// Ambil parameter id dari URL, hasilnya string
	userIDStr := c.Param("id")
	// Ubah string jadi angka (uint64).
	// 10 = basis angka (desimal), 32 = ukuran maksimal bit
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	//Panggil method dari service
	user, err := uc.UserService.GetUserByID(uint(userID))
	// Jika user tidak ditemukan → kirim 404
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
