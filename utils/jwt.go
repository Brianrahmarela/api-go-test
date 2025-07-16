package utils

import (
	"api-go-test/config" //Mengambil config durasi token, secret key, dsb.
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// func menghasilkan JWT token (string), atau error jika gagal.
func GenerateToken(userId uint) (string, error) {
	//claims Adalah data yg dimasukkan ke dalam token JWT, biasanya berupa map (key-value).
	claims := jwt.MapClaims{
		"user_id": userId,                                                   // Menyimpan ID user (agar token bisa dikenali siapa pemiliknya)
		"exp":     time.Now().Add(config.GetJwtExpirationDuration()).Unix(), // .Add(...): method dari time.Time → menambahkan durasi tertentu ke waktu sekarang, krn di env 24jam, Expired (waktu kadaluarsa) → sekarang + durasi token.
		"iat":     time.Now().Unix(),                                        //"iat": Issued at (kapan token dibuat), dalam format Unix timestamp.
		//.Unix(): mengubah waktu hasil penjumlahan itu menjadi angka Unix timestamp (jumlah detik sejak 1 Januari 1970 UTC).
	}
	//jwt.NewWithClaims -> generate new token
	//SigningMethodHS256: Metode enkripsi token dengan algoritma HMAC SHA-256.
	// claims: Data (payload) yang ingin dimasukkan ke token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// config.GetJwtSecret(): membaca/ambil JWT_SECRET_KEY di .env
	fmt.Println("config.GetJwtSecret:", config.GetJwtSecret())
	//.SignedString(...): Token ditandatangani dengan secret key (diambil dari konfigurasi). Hasil akhirnya adalah string token JWT, siap dikirim ke user.
	return token.SignedString(config.GetJwtSecret())
}

// Fungsi ini digunakan untuk memverifikasi token yang dikirim oleh user.
// Input: param tokenString (token dalam bentuk string)
// Output: me return uint = userId dari token, atau error jika tidak valid.
func ValidateToken(tokenString string) (uint, error) {
	// jwt.Parse(): jwt.Parse(...) akan buka/urai token ini jadi bentuk data yang bisa dibaca (bukan string lagi),
	// Memeriksa apakah tanda tangannya benar dgn func(token *jwt.Token), Memeriksa apakah token sudah expired atau belum.
	// Jika semua valid, mengembalikan token (yang sudah dicek), dan err = nil.
	// func(token *jwt.Token): memanggil fungsi untuk minta secret key.  tugasnya kasih secret key ke jwt.Parse.
	// Parameter token di situ diisi otomatis oleh Go, yaitu versi token yang sedang dicek tanda tangannya.
	// parameter token func(token *jwt.Token) adalah isi dari tokenString yang sudah dibuka oleh jwt.Parse (bukan string lagi yg dikirim client, tp hasil urai).
	// (token *jwt.Token): fungsi ini menerima parameter bernama token, bertipe *jwt.Token (pointer ke struct Token dari library JWT).

	// token, err -> jika valid mereturn token yg dicek, jika tdk valid return error,
	// token, err :=  token ini BUKAN string lagi️, Tapi sudah jadi struct *jwt.Token. isinya ada Header, Claims, dll.
	// (interface{}, error): fungsi ini akan mengembalikan dua nilai:
	// interface{} → artinya bisa mengembalikan data bertipe apa saja. digunakan jika tdk tau tipe apa datanya (dalam hal ini: token yang sudah di-parse & dicek).
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// config.GetJwtSecret(): memanggil fungsi dari package config untuk mengambil JWT secret key dan me return secret key.
		// nil: artinya tidak ada error saat memberikan secret key.
		return config.GetJwtSecret(), nil
	})
	//jika tdk valid, return 0 dan nil
	if err != nil {
		return 0, err
	}
	//token.Claims: Ambil isi token, .(jwt.MapClaims): Convert claims menjadi map. token.Valid: Mengecek apakah token benar-benar valid (belum kadaluarsa dan tanda tangan cocok).
	//claims["user_id"].(float64): Karena JWT menyimpan angka sebagai float64, jadi harus dikonversi dulu ke uint.

	//token.Claims.(jwt.MapClaims); ok && token.Valid artinya:
	// token: variabel hasil dari jwt.Parse().
	// .Claims: adalah field dari struct jwt.Token. Isinya adalah payload dari token (yaitu user_id, exp, dll).
	//Tapi: .Claims ini bertipe interface{}, artinya bisa berisi data dalam berbagai tipe — kita harus pastikan isinya benar-benar MapClaims supaya bisa kita pakai.

	// maksudnya:
	// token.Claims → ambil isi token (isinya: user_id, exp, dll).
	// .(jwt.MapClaims) → type assertion untuk memastikan token.Claims bisa dibaca seperti map.
	// ok := ... → boolean hasil apakah type assertion tadi berhasil.
	// ok && token.Valid → cek bahwa dua syarat terpenuhi:
	// type assertion berhasil
	// tokennya valid (tanda tangan cocok, tidak expired, dsb).
	//
	//.(jwt.MapClaims) disebut type assertion. Type assertion adalah cara di Go untuk mengubah tipe data dari interface{} ke tipe yang kamu inginkan, supaya bisa digunakan seperti biasa.
	// Karena kadang di Go, kita menerima data dalam bentuk interface{} (tipe bebas). Tapi untuk dipakai, kita harus tahu tipenya sebenarnya.
	// MapClaims itu bukan tipe biasa, tapi tipe yang didefinisikan dalam package jwt
	// Di JWT, token.Claims bertipe interface{} (kenapa? supaya fleksibel bisa pakai berbagai jenis claims).
	// Tapi kamu tahu kamu simpan token pakai jwt.MapClaims, maka kamu butuh type assertion:
	// Artinya:

	// “Saya yakin nilai ini bertipe jwt.MapClaims, jadi saya mau pakai dia sebagai tipe itu.”
	// Karena token.Claims bertipe interface{}, kita tidak bisa langsung akses key seperti ["user_id"] tanpa mengetahui tipenya. Maka kita lakukan type assertion.

	//ok && token.Valid: ok harus true → artinya konversi type assertion berhasil.
	// token.Valid juga harus true → artinya token: Tanda tangannya benar, Belum expired., Tidak rusak.Jika keduanya benar, maka blok if dijalankan.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// JWT menyimpan data dalam bentuk map[string]interface{}.
		// Artinya, saat ambil key dari map claims["user_id"], Go hanya tahu nilainya itu interface{} (alias: bisa tipe apa saja).
		// Makanya ditulis: claims["user_id"].(float64). Itu disebut type assertion. agar go tau tipenya benar benar float64.
		// Kalau nggak kasih pernyataan  type assertion.(float64)  interface{} tdk bisa langsung diconvert ke uint
		// uint(...) -> Mengubah hasil float64 tadi menjadi uint (unsigned integer), agar sesuai dengan tipe userId di aplikasi.
		// tp pastikan hasilnya floet64 agar bisa diubah ke uint.
		// claims["user_id"], claims: map yang berisi isi dari token (setelah berhasil diverifikasi).
		// ["user_id"]: ambil nilai dari key "user_id".
		//.(float64) -> Ini adalah type assertion lagi. Karena nilai di dalam JWT disimpan dalam tipe float64 (terutama untuk angka), maka user_id di-cast ke float64 dulu.
		userId := uint(claims["user_id"].(float64))
		return userId, nil
	}

	return 0, err
}

// Gunakan type assertion saat:
// ambil data dari interface{} dan mau menggunakannya sebagai tipe tertentu.
// sedang parsing JSON (atau JWT claims) → karena hasilnya selalu interface{}.
// mau mengecek tipe asli data.
