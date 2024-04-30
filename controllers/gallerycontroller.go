package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetGalleryById(c *fiber.Ctx) error {
	// Ambil ID produk dari parameter URL
	productId := c.Params("id")

	// Buat variabel untuk menyimpan galeri yang sesuai dengan ID produk
	var gallery []models.Gallery

	// Cari galeri berdasarkan ID produk
	if err := initialize.DB.Where("product_id = ?", productId).Find(&gallery).Error; err != nil {
		// Jika galeri tidak ditemukan, kirim respons not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusNotFound,
				Status:  "Not Found",
				Message: "Gallery not found for product",
			}
			return c.Status(fiber.StatusNotFound).JSON(response)
		}
		// Jika terjadi kesalahan lain saat mengambil galeri, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to fetch gallery",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Membuat slice untuk menyimpan path galeri
	var paths []string
	for _, gallery := range gallery {
		paths = append(paths, gallery.Path)
	}

	// Mengirimkan respons sukses dengan data path galeri yang ditemukan
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   paths,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
