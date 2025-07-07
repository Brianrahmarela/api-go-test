package middleware

import (
	"api-go-test/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// mengecek token saat user mengakses API
// gin.HandlerFunc: Tipe fungsi khusus dari Gin yang menerima konteks permintaan (*gin.Context)
func AuthMiddleware() gin.HandlerFunc {
	//c *gin.Context: c adalah konteks dari permintaan web (request), berisi informasi seperti header, parameter, dll.
	return func(c *gin.Context) {
		//c.GetHeader(...): Mengambil nilai dari header HTTP bernama "Authorization"
		// authHeader: nilai dari header HTTP. cth: "Bearer eyJhbGciOiJI..."
		authHeader := c.GetHeader("Authorization")
		//jika authHeader kosong, berarti tidak ada token.
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization is required"})
			//c.Abort(): Menghentikan eksekusi; request tidak diteruskan ke handler berikutnya.
			c.Abort()
			//return: Keluar dari middleware.
			return
		}
		// strings.Split(...): Memecah string dengan delimiter spasi.
		// Contoh: "Bearer eyJhbGci..." jadi ["Bearer", "eyJhbGci..."].
		parts := strings.Split(authHeader, " ")
		// len(parts) != 2: Mengecek apakah hasil split adalah array dgn length 2 &
		// parts[0] != "Bearer": Apakah array index pertama bukan kata "Bearer"?
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Authorization header must be Bearer Token"})
			c.Abort()
			return
		}
		// Validasi JWT Token
		// parts[1]: Ambil index kedua dari header -> nilai JWT token.
		// utils.ValidateToken(...): Panggil fungsi untuk memeriksa valid tdknya token.
		// Jika token valid, userId dikembalikan dari token. Jika tidak, err akan berisi error.
		userId, err := utils.ValidateToken(parts[1])
		//Jika token tidak valid atau expired, kirim respons 401.
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}
		// c.Set(...): Menyimpan data ke Context, agar bisa dipakai di handler (misal: controller).
		// "userId": Nama key-nya, bisa diambil nanti dengan c.Get("userId").
		c.Set("userId", userId)
	}
}
