package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateTag(c *fiber.Ctx) error {
	// Parse body request ke dalam struct TagCreateRequest
	var requestBody models.TagRequest
	if err := c.BodyParser(&requestBody); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi request body
	if err := validate.Struct(&requestBody); err != nil {
		errors := make(map[string][]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			message := fmt.Sprintf("%s is %s", field, tag)
			errors[field] = append(errors[field], message)
		}
		response := helpers.ResponseError{
			Code:   fiber.StatusBadRequest,
			Status: "Bad Request",
			Error:  errors,
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Buat entitas Tag dari data request
	tag := models.Tag{
		Tag: requestBody.Tag,
	}

	// Menyimpan tag ke dalam database
	if err := initialize.DB.Create(&tag).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save tag",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Mengirimkan respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Tag created successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteTag(c *fiber.Ctx) error {
	// Ambil ID tag dari parameter URL
	tagID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid tag ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Cari tag berdasarkan ID di database
	var tag models.Tag
	if err := initialize.DB.First(&tag, tagID).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusNotFound,
			Status:  "Not Found",
			Message: "Tag not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Hapus tag dari database
	if err := initialize.DB.Delete(&tag).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to delete tag",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Mengirimkan respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Tag deleted successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
func UpdateTag(c *fiber.Ctx) error {
	// Ambil ID tag dari parameter URL
	tagID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid tag ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Parse body request ke dalam struct TagUpdateRequest
	var requestBody models.TagRequest
	if err := c.BodyParser(&requestBody); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi request body
	if err := validate.Struct(&requestBody); err != nil {
		errors := make(map[string][]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			message := fmt.Sprintf("%s is %s", field, tag)
			errors[field] = append(errors[field], message)
		}
		response := helpers.ResponseError{
			Code:   fiber.StatusBadRequest,
			Status: "Bad Request",
			Error:  errors,
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Cari tag berdasarkan ID di database
	var tag models.Tag
	if err := initialize.DB.First(&tag, tagID).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusNotFound,
			Status:  "Not Found",
			Message: "Tag not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Update tag dengan data baru
	tag.Tag = requestBody.Tag

	// Simpan perubahan ke dalam database
	if err := initialize.DB.Save(&tag).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to update tag",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Mengirimkan respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Tag updated successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
