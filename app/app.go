package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/muchlist/sagasql/config"
	"github.com/muchlist/sagasql/db"
	"github.com/muchlist/sagasql/middle"
	"log"
)

// RunApp menjalankan framework fiber
func RunApp() {

	// Inisiasi database pool
	dbPool := db.InitDB()
	defer dbPool.Close()

	// Inisiasi fiber
	app := fiber.New()

	// memasang middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Content-Type, Accept, Authorization",
	}))

	// file static gambar
	app.Static("/image", "./static/image")

	// url mapping
	api := app.Group("/api/v1")

	//USER
	api.Get("/users/:username", userHandler.Get)
	api.Get("/users", userHandler.Find)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Get("/profile", middle.NormalAuth(), userHandler.GetProfile)
	api.Post("/register-force", userHandler.Register)                                // <- seharusnya gunakan middleware agar hanya admin yang bisa meregistrasi
	api.Post("/register", middle.NormalAuth(config.RoleAdmin), userHandler.Register) // <- hanya admin yang bisa meregistrasi
	api.Put("/users/:username", middle.NormalAuth(config.RoleAdmin), userHandler.Edit)
	api.Delete("/users/:username", middle.NormalAuth(config.RoleAdmin), userHandler.Delete)

	//PRODUCT
	api.Get("/products/:id", middle.NormalAuth(), productHandler.Get)
	api.Get("/products", middle.NormalAuth(), productHandler.Find)
	api.Post("/products", middle.NormalAuth(), productHandler.Insert)
	api.Put("/products/:id", middle.NormalAuth(), productHandler.Edit)
	api.Delete("/products/:id", middle.NormalAuth(), productHandler.Delete)

	if err := app.Listen(":3500"); err != nil {
		log.Fatalf("Aplikasi tidak dapat dijalankan. Error : %s", err.Error())
		return
	}
}
