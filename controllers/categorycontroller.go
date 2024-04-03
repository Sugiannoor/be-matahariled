package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func init() {
	validate = validator.New()
}

func GetCategoriesLabel(c *fiber.Ctx) error {
	// Ambil semua kategori dari database
	var categories []models.Category
	if err := initialize.DB.Find(&categories).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil kategori, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat respons dengan format yang diinginkan
	var categoryOptions []map[string]interface{}
	for _, category := range categories {
		option := map[string]interface{}{
			"value": category.CategoryId,
			"label": category.Category,
		}
		categoryOptions = append(categoryOptions, option)
	}

	// Kembalikan respons sukses dengan data kategori ke klien
	response := helpers.GeneralResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryOptions,
	}
	return c.JSON(response)
}

func GetCountCategory(c *fiber.Ctx) error {
	// Hitung jumlah total kategori dari database
	var count int64
	if err := initialize.DB.Model(&models.Category{}).Count(&count).Error; err != nil {
		// Jika terjadi kesalahan saat menghitung kategori, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kembalikan respons dengan total kategori
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   count,
	}
	return c.JSON(response)
}

func CreateCategory(c *fiber.Ctx) error {
	// Bind request body ke struct Category
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		// Jika terjadi kesalahan dalam memparsing body, kirim respons kesalahan ke klien
		response := helpers.GeneralResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "Kesalahan Format Data",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi data kategori
	if err := validate.Struct(&category); err != nil {
		errors := make(map[string][]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			var tagName string
			switch field {
			case "Category":
				tagName = "category"
			default:
				tagName = field
			} // Mendapatkan nama tag JSON yang sesuai
			message := tagName + " Mohon diisi" // Pesan kesalahan yang disesuaikan
			errors[tagName] = append(errors[field], message)
		}
		response := helpers.ResponseError{
			Code:   400,
			Status: "Bad Request",
			Error:  errors,
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Simpan kategori ke database
	if err := initialize.DB.Create(&category).Error; err != nil {
		// Jika terjadi kesalahan saat menyimpan kategori, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kembalikan respons sukses dengan data kategori yang baru dibuat ke klien
	response := helpers.ResponseMassage{
		Code:    200,
		Status:  "Created",
		Message: "Category Berhasil dibuat",
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

func UpdateCategory(c *fiber.Ctx) error {
	// Bind request body ke struct Category
	var updatedCategory models.Category
	if err := c.BodyParser(&updatedCategory); err != nil {
		// Jika terjadi kesalahan dalam memparsing body, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    400,
			Status:  "Bad Request",
			Message: "Kesalahan Format Pengiriman",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi data kategori yang diperbarui
	if err := validate.Struct(&updatedCategory); err != nil {
		errors := make(map[string][]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			var tagName string
			switch field {
			case "Category":
				tagName = "category"
			default:
				tagName = field
			} // Mendapatkan nama tag JSON yang sesuai
			message := tagName + " Mohon diisi" // Pesan kesalahan yang disesuaikan
			errors[tagName] = append(errors[field], message)
		}
		response := helpers.ResponseError{
			Code:   400,
			Status: "Bad Request",
			Error:  errors,
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Cari kategori berdasarkan ID yang diterima dari body permintaan
	var category models.Category
	if err := initialize.DB.First(&category, updatedCategory.CategoryId).Error; err != nil {
		// Jika kategori tidak ditemukan, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    404,
			Status:  "Not Found",
			Message: "Data tidak tersedia",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Perbarui kategori di database
	if err := initialize.DB.Model(&category).Updates(&updatedCategory).Error; err != nil {
		// Jika terjadi kesalahan saat memperbarui kategori, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kembalikan respons sukses dengan data kategori yang diperbarui ke klien
	response := helpers.ResponseMassage{
		Code:    200,
		Status:  "OK",
		Message: "Category berhasil diupdate",
	}
	return c.JSON(response)
}
