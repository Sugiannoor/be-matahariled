package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllContracts(c *fiber.Ctx) error {
	// Ambil semua data kontrak dari database
	var contracts []models.Contract
	if err := initialize.DB.Preload("Products").Find(&contracts).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to fetch contracts",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Membuat slice untuk menyimpan respons kontrak
	contractResponses := make([]models.ContractResponse, len(contracts))

	// Mengisi data respons kontrak
	for i, contract := range contracts {
		contractResponses[i] = models.ContractResponse{
			ContractId:   contract.ContractId,
			Title:        contract.Title,
			Description:  contract.Description,
			StartDate:    contract.StartDate,
			EndDate:      contract.EndDate,
			ProductNames: models.GetProductNames(contract.Products),
		}
	}
	response := helpers.ResponseGetAll{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   contractResponses,
	}

	// Mengirimkan respons dengan daftar kontrak
	return c.Status(fiber.StatusOK).JSON(response)
}

func GetContractByID(c *fiber.Ctx) error {
	// Ambil ID kontrak dari parameter URL
	contractID, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid contract ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Cari kontrak berdasarkan ID di database
	var contract models.Contract
	if err := initialize.DB.Preload("Products").First(&contract, contractID).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusNotFound,
			Status:  "Not Found",
			Message: "Contract not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Membuat respons berdasarkan kontrak yang ditemukan
	contractResponse := models.ContractResponse{
		ContractId:   contract.ContractId,
		Title:        contract.Title,
		Description:  contract.Description,
		StartDate:    contract.StartDate,
		EndDate:      contract.EndDate,
		ProductNames: models.GetProductNames(contract.Products),
	}

	// Mengirimkan respons dengan data kontrak
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   contractResponse,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
