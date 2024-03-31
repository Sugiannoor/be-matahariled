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
	user.Get("/datatable", controllers.UserDatatable)
	// Group Products
	product := api.Group("/product")
	product.Get("/all", controllers.GetAllProducts)
	product.Get("/", controllers.GetProductById)
	product.Get("/datatable", controllers.GetDatatableProducts)
	product.Post("/", controllers.CreateProduct)
	product.Put("/", controllers.EditProduct)
	product.Delete("/", controllers.DeleteProduct)
	// Group Category
	category := api.Group("/category")
	category.Get("/label", controllers.GetCategoriesLabel)
	category.Post("/", controllers.CreateCategory)
	category.Put("/", controllers.UpdateCategory)

	// Group Contract
	contract := api.Group("/contract")
	contract.Get("/all", controllers.GetAllContracts)

	app.Listen(":8000")

}
