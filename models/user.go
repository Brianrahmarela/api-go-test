package models

import (
	"golang.org/x/crypto/bcrypt" //Digunakan untuk mengenkripsi (hash) dan memverifikasi password dengan aman.
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string   `json:"name"`                                       //json:"name": Data dikirim sebagai JSON via API, field akan disebut "name".
	Email    string   `json:"email" gorm:"unique"`                        //gorm:"unique": Email harus unik (tidak boleh sama dengan user lain di database).
	Password string   `json:"password"`                                   //Kata sandi yang akan disimpan dalam bentuk hash (bukan plaintext).
	Profile  *Profile `json:"profile,omitempty" gorm:"foreignKey:UserID"` //Relasi ke tabel Profile. model User memberitahu bahwa di tabel profiles ada kolom UserID yang merupakan foreign key yang menunjuk ke tabel ini users.id.
}

// Profile Model (one to one relationship with User)
type Profile struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	//binding:"required": Email dan Password wajib diisi saat dikirim lewat API.
}

// Base Request
type CreateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Bio       string `json:"bio"`
}

// Base Response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// (u *User): Fungsi ini milik pointer ke struct User, artinya bisa mengubah isi User langsung.
// password string: Parameter yang menerima password dalam bentuk teks biasa.
// error: Fungsi return error jika terjadi kesalahan.
func (u *User) HashPassword(password string) error {
	//bcrypt.GenerateFromPassword: Fungsi dari library bcrypt untuk membuat hash dari password.
	//[]byte(password): Password diubah jadi array byte.
	//bcrypt.DefaultCost: Tingkat kompleksitas enkripsi standar.
	//hashedPassword: hasil hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	//Password yang sudah di-hash disimpan ke struct User, menggantikan password asli.
	u.Password = string(hashedPassword)
	//Jika semua lancar, kembalikan nil (tidak ada error).
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
