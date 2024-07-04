package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"fmt"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetAllContracts(c *fiber.Ctx) error {
	// Ambil semua data kontrak dari database dengan preloading untuk memuat relasi User dan Products
	var contracts []models.Contract
	if err := initialize.DB.Preload("User").Preload("Products").Find(&contracts).Error; err != nil {
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
		// Mendapatkan nama pengguna dari relasi
		userName := contract.User.FullName // Pastikan atribut 'Name' adalah atribut yang benar untuk menyimpan nama pengguna

		// Mendapatkan daftar nama produk dari relasi
		var productNames []string
		for _, product := range contract.Products {
			productNames = append(productNames, product.Title) // Pastikan atribut 'Name' adalah atribut yang benar untuk menyimpan nama produk
		}

		contractResponses[i] = models.ContractResponse{
			ContractId:   contract.ContractId,
			Title:        contract.Title,
			Description:  contract.Description,
			StartDate:    contract.StartDate,
			EndDate:      contract.EndDate,
			UserName:     userName,
			ProductNames: productNames,
		}
	}

	// Mengirimkan respons dengan daftar kontrak beserta nama pengguna dan daftar produk
	response := helpers.ResponseGetAll{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   contractResponses,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetCountContract(c *fiber.Ctx) error {
	// Hitung jumlah total kontrak dari database
	var count int64
	if err := initialize.DB.Model(&models.Contract{}).Count(&count).Error; err != nil {
		// Jika terjadi kesalahan saat menghitung kontrak, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kembalikan respons dengan total kontrak
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   count,
	}
	return c.JSON(response)
}

func CreateContract(c *fiber.Ctx) error {
	// Parse body request ke dalam struct ContractCreateRequest
	var requestBody models.ContractRequest
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

	if len(requestBody.ProductIDs) == 0 {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Product id setidaknya ada 1",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Membuat entitas Contract dari data request
	contract := models.Contract{
		Title:       requestBody.Title,
		Description: requestBody.Description,
		StartDate:   requestBody.StartDate,
		EndDate:     requestBody.EndDate,
		UserID:      requestBody.UserID,
	}

	if requestBody.UserID != 0 {
		var user models.User
		if err := initialize.DB.First(&user, requestBody.UserID).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusBadRequest,
				Status:  "Bad Request",
				Message: fmt.Sprintf("User with ID %d not found", requestBody.UserID),
			}
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		contract.UserID = requestBody.UserID
	}

	// Membuat relasi dengan produk (jika ada)
	for _, productID := range requestBody.ProductIDs {
		var product models.Product
		if err := initialize.DB.First(&product, productID).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusBadRequest,
				Status:  "Bad Request",
				Message: fmt.Sprintf("Product with ID %d not found", productID),
			}
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		contract.Products = append(contract.Products, product)
	}

	// Menyimpan kontrak ke dalam database
	if err := initialize.DB.Create(&contract).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Mengirimkan respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Contract created successfully",
	}
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

func GetContractsDataTable(c *fiber.Ctx) error {
	// Ambil nilai parameter limit, sort, sort_by, dan search dari query string
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	search := c.Query("search")
	userID, _ := strconv.Atoi(c.Query("user_id"))

	// Tentukan default nilai jika parameter tidak ada
	if limit <= 0 {
		limit = 10 // Nilai default untuk limit adalah 10
	}
	offset := (page - 1) * limit
	// Lakukan pengambilan data dari database dengan menggunakan parameter limit, sort, dan sort_by
	var contracts []models.Contract
	query := initialize.DB.Preload("User").Preload("Products")

	// Jika parameter userID tidak kosong, tambahkan filter berdasarkan userID
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	// Jika parameter search tidak kosong, tambahkan filter pencarian berdasarkan judul kontrak, nama pengguna, dan nama produk
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Jika parameter sort dan sort_by disediakan, lakukan pengurutan berdasarkan kolom yang dimaksud
	if sort != "" && sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sort))
	}

	var totalRecords int64
	if err := initialize.DB.Model(&models.Contract{}).Count(&totalRecords).Error; err != nil {
		response := helpers.GeneralResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   nil,
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	// Limit jumlah data yang diambil sesuai dengan nilai parameter limit
	query = query.Limit(limit).Offset(offset)

	// Lakukan pengambilan data
	if err := query.Find(&contracts).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil data kontrak, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	// Membuat slice untuk menyimpan respons data kontrak
	contractResponses := make([]models.ContractResponseDatatable, len(contracts))

	// Mengisi data respons kontrak
	for i, contract := range contracts {
		// Mendapatkan nama pengguna dari relasi
		userName := contract.User.FullName // Pastikan atribut 'Name' adalah atribut yang benar untuk menyimpan nama pengguna

		// Mendapatkan daftar nama produk dari relasi
		var productNames []map[string]interface{}
		for _, product := range contract.Products {
			productInfo := map[string]interface{}{
				"name": product.Title,
				"id":   product.ProductId,
			}
			productNames = append(productNames, productInfo)
		}

		contractResponses[i] = models.ContractResponseDatatable{
			ContractId:   contract.ContractId,
			Title:        contract.Title,
			Description:  contract.Description,
			StartDate:    contract.StartDate,
			EndDate:      contract.EndDate,
			UserId:       contract.UserID,
			UserName:     userName,
			ProductNames: productNames,
		}
	}
	data := make([]interface{}, len(contractResponses))
	for i, v := range contractResponses {
		data[i] = v
	}
	// Kembalikan respons datatable
	response := helpers.DataTableResponse{
		CurrentPage:  1,              // Nomor halaman saat ini (default 1)
		FirstPageURL: "",             // URL halaman pertama
		From:         1,              // Nomor record pertama pada halaman saat ini
		LastPage:     totalPages,     // Total jumlah halaman (default 1)
		LastPageURL:  "",             // URL halaman terakhir
		NextPageURL:  "",             // URL halaman berikutnya
		PrevPageURL:  "",             // URL halaman sebelumnya
		To:           len(contracts), // Nomor record terakhir pada halaman saat ini
		Total:        len(contracts), // Total jumlah record
		Data:         data,           // Data kontrak
	}

	// Kembalikan respons umum dengan data datatable
	return c.JSON(helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func UpdateContract(c *fiber.Ctx) error {
	// Ambil id kontrak dari parameter URL
	contractID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid contract ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Parse body request ke dalam struct ContractRequest
	var requestBody models.ContractRequest
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

	// Memeriksa apakah kontrak dengan ID yang diberikan ada di database
	var existingContract models.Contract
	if err := initialize.DB.Preload("Products").First(&existingContract, contractID).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusNotFound,
			Status:  "Not Found",
			Message: "Contract not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Update data kontrak dengan data baru
	existingContract.Title = requestBody.Title
	existingContract.Description = requestBody.Description
	existingContract.StartDate = requestBody.StartDate
	existingContract.EndDate = requestBody.EndDate

	// Jika ada UserID yang disertakan, validasi apakah UserID tersebut ada
	if requestBody.UserID != 0 {
		var user models.User
		if err := initialize.DB.First(&user, requestBody.UserID).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusBadRequest,
				Status:  "Bad Request",
				Message: fmt.Sprintf("User with ID %d not found", requestBody.UserID),
			}
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		existingContract.UserID = requestBody.UserID
	}

	// Menghapus semua relasi produk yang terkait dengan kontrak
	if err := initialize.DB.Model(&existingContract).Association("Products").Clear(); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to update contract",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Menambahkan relasi dengan produk baru (jika ada)
	for _, productID := range requestBody.ProductIDs {
		var product models.Product
		if err := initialize.DB.First(&product, productID).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusBadRequest,
				Status:  "Bad Request",
				Message: fmt.Sprintf("Product with ID %d not found", productID),
			}
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		existingContract.Products = append(existingContract.Products, product)
	}

	// Menyimpan perubahan kontrak ke dalam database
	if err := initialize.DB.Save(&existingContract).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to update contract",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Mengirimkan respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Contract updated successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteContract(c *fiber.Ctx) error {
	// Ambil parameter ID kontrak dari URL
	contractID := c.Params("id")

	// Mulai transaksi
	tx := initialize.DB.Begin()

	// Hapus hubungan antara kontrak dan produk di tabel ContractProduct
	if err := tx.Where("contract_contract_id = ?", contractID).Delete(&models.ContractProduct{}).Error; err != nil {
		tx.Rollback()
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to delete contract's products",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Hapus kontrak
	if err := tx.Where("contract_id = ?", contractID).Delete(&models.Contract{}).Error; err != nil {
		tx.Rollback()
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to delete contract",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Commit transaksi
	tx.Commit()

	// Mengirimkan respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Contract deleted successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
