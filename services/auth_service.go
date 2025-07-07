// Package ini berisi logika bisnis utama, seperti proses register dan login.
package services

import (
	"api-go-test/models" //models: Mengimpor model User dan LoginRequest.
	"api-go-test/utils"  //utils: Digunakan untuk memanggil GenerateToken (membuat JWT).
	"errors"             //errors: Package bawaan Go untuk membuat error baru.

	"gorm.io/gorm" //gorm.io/gorm: ORM untuk berinteraksi dengan database.
)

// "wadah" logika login & register dengan akses ke database.
type AuthService struct {
	DB *gorm.DB // Menyimpan koneksi database (*gorm.DB) di dalam AuthService.
}

// Fungsi untuk membuat instance baru dari AuthService.
func NewAuthService(db *gorm.DB) *AuthService {
	// Mengembalikan pointer AuthService dengan koneksi DB disimpan di dalamnya.
	return &AuthService{DB: db}
}

// membuat method milik struct AuthService.
// user *models.User: Parameter input data user.
// returns (string, error): Hasilnya adalah JWT token dan kemungkinan error.
func (as *AuthService) Register(user *models.User) (string, error) {
	// Mengubah password asli menjadi hash. Jika gagal hashing, langsung return error.
	if err := user.HashPassword(user.Password); err != nil {
		return "", errors.New("error hashing password")
	}

	// Create user in database
	if err := as.DB.Create(user).Error; err != nil {
		// jika gagal insert, return error.
		return "", errors.New("error creating user")
	}

	// utils.GenerateToken: Membuat token JWT berdasarkan user.ID.
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", errors.New("error generating token")
	}
	//Jika berhasil, kembalikan token dan tidak ada error (nil).
	return token, nil
}

// Fungsi untuk memproses login user berdasarkan email & password.
// loginReq: struct yang berisi input login dari us
func (as *AuthService) Login(loginReq *models.LoginRequest) (string, error) {
	var user models.User

	// Mencari user berdasarkan email di database.
	// if err :=: jika ada error, Buat variabel err dan isi dengan hasil error dari proses pencarian user di database.
	// as.DB Adalah objek koneksi database (dari GORM), disimpan di struct AuthService.
	// .Where("email = ?", loginReq.Email): bikin query SQL WHERE email = ?.
	// ? ini nilainya diambil dari input user.
	// .First(&user): .First fungsi dari gorm utk ambil baris pertama dari hasil pencarian,
	// dan simpan otomatis ke struct user.
	// .Error: Properti dari GORM yang mengambil hasil error dari query tadi.
	// .Error dipakai karena: Data sudah disimpan ke &user, tapi hanya kalau tidak error.
	// GORM tidak akan memberitahu secara otomatis apakah query berhasil atau gagal.
	// harus mengeceknya sendiri dengan melihat .Error nya
	// Jika tidak ada error/query berhadil → nil, Jika user tidak ditemukan → gorm.ErrRecordNotFound,
	// Jika koneksi DB gagal → error lain
	// ;: Akhir dari pernyataan if err := ....
	//err != nil: Mengecek apakah ada error. Jika err tidak kosong, maka proses selanjutnya adalah:
	//if errors.Is(err, gorm.ErrRecordNotFound) {
	if err := as.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//Jika tidak ketemu, return 2 yaitu string kosong "" & "invalid email or password".
			// sesuai return fungsi Login (string, error)
			return "", errors.New("invalid email or password")
		}
		return "", err
	}
	// contoh pengecekan sederhana if diatas yg lebih bisa dibaca:
	// 1. Jalankan query, simpan jika ada ke &user dan tangkap .Error,
	// Dalam GORM, setiap query seperti .First(), .Create(), .Where() mengembalikan struct bernama *gorm.DB
	// yg punya banyak properti, salah satunya: .Error yg lgsg disimpan ke variable err
	// err := db.Where("email = ?", loginReq.Email).First(&user).Error

	// 2. Periksa error-nya
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		fmt.Println("User tidak ditemukan")
	// 	} else {
	// 		fmt.Println("Terjadi error:", err)
	// 	}
	// }

	// Check password, Cocokkan password input dengan password hash di database.
	if err := user.CheckPassword(loginReq.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate token, Jika password cocok, buat token JWT.
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", errors.New("error generating token")
	}

	return token, nil
}
