package services

import (
	"api-go-test/models"
	"errors"

	"gorm.io/gorm"
)

// struct { DB *gorm.DB } → struct ini punya satu field:
type UserService struct {
	DB *gorm.DB //pointer ke gorm.DB, menyimpan koneksi ke database dari GORM.
}

// NewUserService → constructor (pembuat objek struct)
// (db *gorm.DB) → menerima parameter berupa pointer ke koneksi database.
// *UserService → Tipe data return-nya: yaitu tipe data pointer ke struct UserService.
// fungsi ini hanya return tipe data pointer ke struct UserService, bukan nilai struct-nya langsung/ alamat.
// agar bisa diakses dari mana pun. misalnya ke controller.
func NewUserService(db *gorm.DB) *UserService {
	// membuat struct baru UserService isi DB-nya dengan db, dan kembalikan alamat memori-nya ke pemanggil fungsi NewUserService.
	// membuat struct baru UserService dgnUserService{DB: db},
	// UserService adalah nama struct-nya.
	// {DB: db} adalah pengisian nilai field dalam struct-nya, hasil cth:
	// UserService{DB: 0xc000123abc} -> masih berbentuk nilai biasa (bukan pointer).
	// gunakan &UserService utk Ambil alamat memori dari struct yang baru dibuat
	// tujuannya: Tidak copy struct-nya, Bisa digunakan di tempat lain (fungsi, controller, dll)
	return &UserService{DB: db}
	// cth pemakaian: userSvc := NewUserService(db)
	// NewUserService(db) mengembalikan alamat struct UserService
	// userSvc adalah variabel bertipe *UserService (pointer)
	// Jadi userSvc menyimpan alamat memori dari UserService yang dibuat
	// krn userSvc menyimpan alamat, bisa diubah isinya:
	// cth: userSvc.DB = nil
	//tidak perlu menulis *userSvc.DB = nil?
	//krn Go melakukan dereference otomatis saat akses field dari struct pointer.
	//bisa jg (*userSvc).DB = nil  // bentuk eksplisit (panjang)
}

// func (us *UserService) → method milik objek UserService (dengan alias us).
// GetAllUsers() → nama method-nya.
// ([]models.User, error) → return-nya adalah:
// slice of User dari package models & error jika ada kesalahan
func (us *UserService) GetAllUsers() ([]models.User, error) {
	// mendeklarasikan slice kosong users yang akan diisi dari DB.
	var users []models.User
	// us.DB.Find(&users) → ambil semua data user dari database.
	// .Error → ambil error dari operasi Find().
	// if err != nil → kalau ada error, kembalikan nil dan error tersebut.
	if err := us.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	// Remove passwords from all users for security
	// for i := range users → loop setiap user.
	for i := range users {
		//kosongkan field password untuk keamanan.
		users[i].Password = ""
	}
	// Kembalikan data user (tanpa password), dan error nil artinya sukses.
	return users, nil
}

func (us *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	// First(&user, userID) → ambil satu record berdasarkan ID.
	// .Error → ambil error dari hasil query.
	if err := us.DB.First(&user, userID).Error; err != nil {
		//errors.Is(...) → cek apakah error-nya karena tidak ditemukan.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Kalau iya, kembalikan error "user not found".
			return nil, errors.New("user not found")
		}
		//Kalau bukan error itu, kembalikan error apa pun yang terjadi.
		return nil, err
	}

	// Kosongkan password sebelum dikirim agar aman
	user.Password = ""
	//Kembalikan pointer ke user, dan nil sebagai error (berarti sukses).
	return &user, nil
}
