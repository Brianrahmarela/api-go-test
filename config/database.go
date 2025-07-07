package config

import (
	"fmt" //Untuk mencetak (print) ke terminal atau membuat string.
	"os"  //Untuk ambil environment variable (pakai os.Getenv).

	"github.com/joho/godotenv" //Package utk baca env
	"gorm.io/driver/mysql"     //Driver MySQL untuk GORM.
	"gorm.io/gorm"             //Package GORM (ORM untuk Go).
)

// agar bisa membaca konfigurasi gorm db, *gorm.DB: Mengembalikan pointer ke objek database GORM.
func ConnectDatabase() *gorm.DB {
	// godotenv utk baca file .env & memasukkan isinya ke environment (agar os.Getenv bisa membacanya).
	// Load() bukan utk mengembalikan data, tapi untuk melakukan aksi, dan lapor error jika gagal.
	// Kalau berhasil, dia langsung masukkan ke os.Getenv tanpa return nilai.
	errEnv := godotenv.Load()
	if errEnv != nil {
		fmt.Println("Gagal load .env:", errEnv)
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	//dsn = Data Source Name, format string koneksi ke MySQL yg dibutuhkan driver MySQL untuk membuka koneksi
	//fmt.Sprintf(data dari env): Gabung string berdasarkan format.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)
	// hasil print dsn: ("%s:%s"= dbUser, dbPass, @tcp(%s:%s)=dbHost, dbPort, /%s?=dbName)
	// root:1234@tcp(localhost:3306)/crud_jwt_go?charset=utf8mb4&parseTime=True&loc=Local
	// format yg lebih mudah dibaca: dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/"
	// + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println("DSN (Data source name) yang digunakan:", dsn)
	// gorm.Open(...): Fungsi untuk membuka koneksi database.
	// mysql.Open(dsn): Gunakan driver MySQL dan string koneksi dsn.
	// &gorm.Config{}: Gunakan konfigurasi default dari GORM.
	// fungsi gorm.Open(...) mengembalikan 2 nilai (db, err):
	// db → variabel untuk menyimpan objek koneksi GORM.
	// err → variabel untuk menangkap error kalau koneksi gagal.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// Jika gagal konek, tampilkan pesan error dan hentikan program (panic).
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	// jika berhasil Kembalikan objek koneksi database ke pemanggil fungsi.
	return db
}
