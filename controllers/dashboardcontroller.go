package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"

	"github.com/gofiber/fiber/v2"
)

func GetDashboard(c *fiber.Ctx) error {
	// Hitung jumlah total produk dari database
	var productCount int64
	if err := initialize.DB.Model(&models.Product{}).Count(&productCount).Error; err != nil {
		// Jika terjadi kesalahan saat menghitung produk, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Hitung jumlah total pengguna dari database
	var userCount int64
	if err := initialize.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		// Jika terjadi kesalahan saat menghitung pengguna, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Hitung jumlah total kontrak dari database
	var contractCount int64
	if err := initialize.DB.Model(&models.Contract{}).Count(&contractCount).Error; err != nil {
		// Jika terjadi kesalahan saat menghitung kontrak, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Hitung jumlah total riwayat dari database
	var historyCount int64
	if err := initialize.DB.Model(&models.History{}).Count(&historyCount).Error; err != nil {
		// Jika terjadi kesalahan saat menghitung riwayat, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Hitung jumlah total kategori dari database
	var categoryCount int64
	if err := initialize.DB.Model(&models.Category{}).Count(&categoryCount).Error; err != nil {
		// Jika terjadi kesalahan saat menghitung kategori, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kembalikan respons dengan data dashboard
	dashboard := map[string]int64{
		"product_count":  productCount,
		"user_count":     userCount,
		"contract_count": contractCount,
		"history_count":  historyCount,
		"category_count": categoryCount,
	}
	initializeresponse := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   dashboard,
	}
	return c.JSON(initializeresponse)
}
