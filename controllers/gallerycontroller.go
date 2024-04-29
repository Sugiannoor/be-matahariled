package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"

	"github.com/gofiber/fiber/v2"
)

func GetGalleriesByProductID(c *fiber.Ctx) error {
	// Ambil ID produk dari parameter URL
	productID := c.Params("id")

	// Cari galeri berdasarkan ID produk
	var galleries []models.Gallery
	if err := initialize.DB.Where("product_id = ?", productID).Find(&galleries).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to fetch galleries",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kirim respons sukses dengan galeri yang ditemukan
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   galleries,
	}
	return c.JSON(response)
}