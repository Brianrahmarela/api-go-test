package routes

import (
	"api-go-test/controllers"
	"api-go-test/middleware"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Membuat instance (object) dari struct UserController. yaitu membuat object yang nanti punya method-method seperti GetUsers, GetUserByID, dll.
	// bentuknya:
	// &UserController{
	// 	UserService: &UserService{
	// 		DB: db,
	// 	},
	// }
	// NewUserController adalah semacam constructor/helper untuk mengisi field (seperti UserService) dalam UserController.
	userController := controllers.NewUserController(db)
	profileController := controllers.NewProfileController(db)

	// router.Group("/"): Fungsi ini membuat grup route baru dengan prefix "/". Semua route di dalam grup ini akan memiliki awalan /.
	// protected :=: Variabel protected menyimpan grup route tersebut. Grup route ini akan digunakan untuk route yang memerlukan proteksi atau middleware (seperti autentikasi).
	protected := router.Group("/")
	// router.Group("/"): grup ini tidak memiliki middleware atau proteksi yang diterapkan.
	// nonProtected :=: Variabel nonProtected menyimpan grup route ini. Grup route ini akan digunakan untuk route yang tidak memerlukan autentikasi atau proteksi lainnya.
	nonProtected := router.Group("/")
	// protected.Use(middleware.AuthMiddleware()): Fungsi ini akan menerapkan middleware AuthMiddleware() pada grup route protected.
	// middleware.AuthMiddleware(): Middleware ini berfungsi untuk memeriksa apakah pengguna sudah terautentikasi atau belum, sebelum route di dalam grup ini dijalankan. Jika pengguna tidak terautentikasi, maka permintaan akan diblokir.
	// Dengan kata lain, semua route yang ada di dalam protected akan membutuhkan autentikasi terlebih dahulu.
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", userController.GetUsers)
		protected.POST("/users/:id/profile", profileController.CreateProfile)
		protected.GET("/users/:id/profile", profileController.GetProfile)
	}

	nonProtected.GET("/users/:id", userController.GetUserByID)
}
