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
	app.Static("/public", "./public")

	// Dashboard
	api.Get("/dashboard", controllers.GetDashboard)

	// Group Auth
	auth := api.Group("/auth")
	auth.Get("/profile", controllers.Index)
	auth.Post("/register", controllers.CreateUser)

	// User
	user := api.Group("/user")
	user.Delete("/", controllers.DeleteUser)
	user.Get("/all", controllers.Index)
	user.Get("/count", controllers.GetCountUser)
	user.Put("/", controllers.EditUser)
	user.Get("/", controllers.GetUserById)
	user.Get("/datatable", controllers.UserDatatable)
	user.Get("/label", controllers.GetUsersLabel)
	// Group Products
	product := api.Group("/product")
	product.Get("/all", controllers.GetAllProducts)
	product.Get("/", controllers.GetProductById)
	product.Get("/count", controllers.GetCountProduct)
	product.Get("/datatable", controllers.GetDatatableProducts)
	product.Get("/label", controllers.GetProductsLabel)
	product.Post("/", controllers.CreateProduct)
	product.Put("/:id", controllers.UpdateProduct)
	product.Delete("/", controllers.DeleteProduct)
	// Group Category
	category := api.Group("/category")
	category.Get("/label", controllers.GetCategoriesLabel)
	category.Get("/count", controllers.GetCountCategory)
	category.Post("/", controllers.CreateCategory)
	category.Put("/", controllers.UpdateCategory)

	// Group Contract
	contract := api.Group("/contract")
	contract.Get("/all", controllers.GetAllContracts)
	contract.Get("/count", controllers.GetCountContract)
	contract.Get("/datatable", controllers.GetContractsDataTable)
	contract.Post("/", controllers.CreateContract)
	contract.Put("/:id", controllers.UpdateContract)
	contract.Delete("/:id", controllers.DeleteContract)

	//Group History
	history := api.Group("/history")
	history.Get("/all", controllers.GetAllHistories)
	history.Get("/count", controllers.GetCountHistory)
	history.Get("/datatable", controllers.GetDatatableHistories)
	history.Post("/", controllers.CreateHistory)
	history.Put("/:id", controllers.UpdateHistory)
	history.Delete("/:id", controllers.DeleteHistory)

	app.Listen(":8000")

}
