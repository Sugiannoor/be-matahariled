package main

import (
	"Matahariled/controllers"
	"Matahariled/initialize"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	initialize.ConnectDatabase()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
	}))
	// Base
	api := app.Group("/api")

	// Group Auth
	auth := api.Group("/auth")
	auth.Get("/profile", controllers.Index)
	auth.Post("/register", controllers.CreateUser)

	// User
	user := api.Group("/user")
	user.Delete("/", controllers.DeleteUser)
	user.Get("/all", controllers.Index)
	user.Put("/", controllers.EditUser)
	user.Get("/", controllers.GetUserById)
	// Group Products

	app.Listen(":8000")
}
