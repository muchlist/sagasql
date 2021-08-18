package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

// RunApp menjalankan framework fiber
func RunApp() {

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

	// todo url mapping

	if err := app.Listen(":3500"); err != nil {
		log.Fatalf("Aplikasi tidak dapat dijalankan. Error : %s", err.Error())
		return
	}
}
