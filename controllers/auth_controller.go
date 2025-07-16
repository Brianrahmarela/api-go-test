package controllers

import (
	"api-go-test/models"
	"api-go-test/services" // Mengimpor service, yang menangani logika login/register.
	"fmt"
	"net/http" // Package standar go untuk kode HTTP (kode status seperti 200, 400, dll).

	"github.com/gin-gonic/gin" //Framework web bernama Gin. Memudahkan buat REST API.
	"gorm.io/gorm"             //ORM untuk database. GORM mempermudah kerja dengan database seperti MySQL.
)

// INI HANYA CETAKAN (BLUE PRINT), menjelaskan apa field-nya, tapi belum berisi nilai.
// type: Kata kunci untuk membuat tipe data struct.
// stuct itu Mendefinisikan “bentuk” data yang punya satu atau lebih field (properti).
// AuthController: Nama tipe/struct yg menyimpan dependensi AuthService. Dugunakan untuk autentikasi (login/register).
// Controller bertanggung jawab pada HTTP layer (menerima permintaan, mengirim respons),
type AuthController struct {
	// Ini menggunakan satu struct dari services di dalam struct lain yaitu AuthController.
	// AuthService: field di dlm controler struct AuthController.
	// *services.AuthService: tipenya adalah Objek pointer dari AuthService (logika login/register).
	// Service bertanggung jawab pada logika bisnis (hash password, database query)
	// & method di controller bisa langsung memanggil fungsi-fungsi logika bisnis itu.
	AuthService *services.AuthService
}

// INI DIBUAT INSTANCENYA & ISI NILAI dgn constructor function (fungsi pembuat instance),
// instance yaitu -> “nilai sbnrnya yg dibuat dari blueprint AuthController dgn field2 yg udh terisi.
//
//	Tujuannya:
//
// -Menyuntikkan (inject) dependensi: Panggil NewAuthService(db) → dapat *AuthService dengan DB terisi ke dalam AuthService di dalam controller.
// -Menyederhanakan inisialisasi: setiap kali kita butuh AuthController, cukup panggil NewAuthController(db),
// tidak perlu menulis ulang kode pembuatan AuthService.
func NewAuthController(db *gorm.DB) *AuthController {
	//return &AuthController{...}: Mengembalikan objek baru AuthController, isinya:
	return &AuthController{
		// AuthService: Dibuat dari services.NewAuthService(db).
		//AuthService juga struct, belum punya nilai DB sampai dipanggil NewAuthService(db)
		AuthService: services.NewAuthService(db),
	}
}

// INI MEMBUAT METHOD (fungsi dengan receiver) PUNYANYA STRUCT AuthController.
// ac *AuthController adalah receiver, mirip this di JS. Artinya, memanggil authController.Register(c) untuk menjalankan logic register.
// Register: Nama fungsi. Untuk registrasi user baru.
// (c *gin.Context): c adalah objek context dari Gin — mewakili request & response.
func (ac *AuthController) Register(c *gin.Context) {
	// Buat objek kosong User dari model yg udh dibuat di folder models
	var user models.User

	// CEK BODY DARI CLIENT APAKAH JSON dgn c.ShouldBindJSON
	// c.ShouldBindJSON(&user): Baca body JSON dari request, masukkan ke struct user.
	if err := c.ShouldBindJSON(&user); err != nil {
		//err != nil: Kalau gagal (misalnya format JSON salah), kirim error HTTP 400.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ac.AuthService.Register(...): Panggil fungsi Register dari service.
	//jika berhasil simpan di variable token
	token, status, err := ac.AuthService.Register(&user)
	fmt.Println("status", status)
	fmt.Println("err AuthService.Register", err)
	if err != nil {
		c.JSON(status, gin.H{
			"status": status,
			"error":  err.Error(),
		})
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//http.StatusCreated: HTTP status 201 (berhasil dibuat).
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginReq models.LoginRequest

	// harus bentuknya JSON
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use service to login user
	// testing123
	token, err := ac.AuthService.Login(&loginReq)
	if err != nil {
		if err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
