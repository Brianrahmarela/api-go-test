package main

import (
	"api-go-test/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()
	fmt.Println("Koneksi database berhasil:", db != nil)
	// gin.Default() → membuat instance router dari Gin. otomatis mengaktifkan logging dan recovery (anti crash).
	// r → variabel router yang akan menangani semua request.
	r := gin.Default()
	//c *gin.Context → parameter yang diberikan oleh Gin untuk tiap request.
	//*gin.Context adalah objek yang berisi semua informasi tentang request & response.
	// Dengan c (alias dari *gin.Context), yg dpt dilakukan:
	// Fungsi	            Keterangan
	// c.Param("id")	    Ambil parameter dari URL
	// c.Query("q")	        Ambil query string dari URL
	// c.PostForm("name")	Ambil data dari form POST
	// c.BindJSON(...)	    Ambil JSON dari request body
	// c.JSON(...)	        Kirim response JSON ke client
	// c.String(...)	    Kirim response string biasa
	// c.Status(...)	    Set status code aja

	// Jadi Context ini seperti wadah komunikasi antara request & response di dalam Gin.
	r.GET("/", func(c *gin.Context) {
		// gin.H → alias untuk map[string]interface{}, artinya { "message": "..." } dalam bentuk Go.
		c.JSON(200, gin.H{"message": "Berhasil get ke / router"})
	})
	r.Run(":8080") // server akan berjalan di port 8080
}
