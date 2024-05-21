package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetHero(c *fiber.Ctx) error {
	var heroes []models.Hero

	// Mengambil data hero dengan limit 5
	if err := initialize.DB.Select("path", "product_id").Limit(5).Find(&heroes).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil data, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to retrieve heroes",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat respons dengan format yang diinginkan
	var heroOptions []map[string]interface{}
	for _, hero := range heroes {
		option := map[string]interface{}{
			"path":       hero.Path,
			"product_id": hero.ProductId,
		}
		heroOptions = append(heroOptions, option)
	}

	// Kembalikan respons sukses dengan data hero ke klien
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   heroOptions,
	}
	return c.JSON(response)
}

func CreateHero(c *fiber.Ctx) error {
	// Parse inputan dari form
	var requestBody struct {
		ProductId int64 `form:"product_id"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	// Validasi apakah ID produk valid
	var product models.Product
	if err := initialize.DB.First(&product, requestBody.ProductId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	// Periksa apakah file diunggah
	file, err := c.FormFile("file")
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "File is required",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Simpan file yang diunggah ke folder public
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	filePath := fmt.Sprintf("/public/%s", filename)
	if err := c.SaveFile(file, fmt.Sprintf("./public/%s", filename)); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save file",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat entitas Hero untuk disimpan dalam database
	hero := models.Hero{
		Path:      filePath,
		Hero_name: filename,
		Size:      strconv.FormatInt(file.Size, 10),
		Format:    filepath.Ext(file.Filename),
		ProductId: requestBody.ProductId,
	}

	// Simpan data hero ke dalam database
	if err := initialize.DB.Create(&hero).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save hero data",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kirim respons sukses
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   hero,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
