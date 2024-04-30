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
	productId := c.Params("id")
	var gallery []models.Gallery

	// Cari galeri berdasarkan ID produk
	if err := initialize.DB.Where("product_id = ?", productId).Find(&gallery).Error; err != nil {
		// Jika galeri tidak ditemukan, kirim respons not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Mengirimkan respons sukses dengan data galeri kosong
			response := helpers.GeneralResponse{
				Code:   fiber.StatusOK,
				Status: "OK",
				Data:   []string{},
			}
			return c.Status(fiber.StatusOK).JSON(response)
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
	if len(paths) == 0 {
		paths = []string{}
	}
	// Mengirimkan respons sukses dengan data path galeri yang ditemukan
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   paths,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
