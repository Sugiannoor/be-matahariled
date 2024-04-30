package main

import (
	"Matahariled/controllers"
	"Matahariled/initialize"
	"Matahariled/middleware"

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
	auth.Get("/profile", controllers.GetProfileHandler)
	auth.Post("/login", controllers.LoginHandler)
	auth.Post("/register", controllers.CreateUser)

	// User
	user := api.Group("/user")
	user.Delete("/", controllers.DeleteUserT)
	user.Get("/all", middleware.JWTMiddleware("Admin", "SuperAdmin"), controllers.Index)
	user.Get("/count", controllers.GetCountUser)
	user.Put("/", controllers.EditUser)
	user.Get("/", controllers.GetUserById)
	user.Get("/datatable", controllers.UserDatatable)
	user.Get("/label", controllers.GetUsersLabel)
	user.Post("/", controllers.CreateUserForm)
	// Group Products
	product := api.Group("/product")
	product.Get("/all", controllers.GetAllProducts)
	product.Get("/count", controllers.GetCountProduct)
	product.Get("/datatable", controllers.GetDatatableProducts)
	product.Get("/label", controllers.GetProductsLabel)
	product.Post("/", controllers.CreateProductT)
	product.Put("/:id", controllers.UpdateProduct)
	product.Delete("/", controllers.DeleteProductT)
	product.Get("/:id", controllers.GetProductById)
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
	history.Get("/user", controllers.GetAllUserPortfolios)
	history.Get("/:id", controllers.GetHistoryById)
	history.Get("/product/:id", controllers.GetHistoryByIdProduct)
	history.Post("/", controllers.CreateHistory)
	history.Put("/:id", controllers.UpdateHistory)
	history.Delete("/:id", controllers.DeleteHistory)

	tag := api.Group("/tag")
	tag.Post("/", controllers.CreateTag)
	tag.Get("/", controllers.GetAllTag)
	tag.Get("/label", controllers.GetTagLabel)

	video := api.Group("/video")
	video.Get("/all", controllers.GetAllVideos)
	video.Get("/datatable", controllers.GetDatatableVideos)

	gallery := api.Group("/gallery")
	gallery.Get("/:id", controllers.GetGalleryById)

	app.Listen(":8000")

}
