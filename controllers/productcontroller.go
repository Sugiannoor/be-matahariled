package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllProducts(c *fiber.Ctx) error {
	// Ambil semua produk dari database
	var products []models.Product
	if err := initialize.DB.Preload("Category").Preload("File").Find(&products).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil produk, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	var customProducts []models.ProductResponse
	for _, product := range products {
		customProduct := models.ProductResponse{
			ProductId:   product.ProductId,
			Title:       product.Title,
			Description: product.Description,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			FileId:      product.FileId,
			CategoryId:  product.CategoryId,
			PathFile:    product.File.Path,
			Category:    product.Category.Category,
		}
		customProducts = append(customProducts, customProduct)
	}

	// Jika tidak ada kesalahan, kirim respons sukses dengan produk yang ditemukan ke klien
	response := helpers.GeneralResponse{
		Code:   200,
		Status: "OK",
		Data:   customProducts,
	}
	return c.JSON(response)
}

func GetDatatableProducts(c *fiber.Ctx) error {
	// Ambil nilai parameter limit, sort, sort_by, dan search dari query string
	limit, _ := strconv.Atoi(c.Query("limit"))
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	search := c.Query("search")

	// Tentukan default nilai jika parameter tidak ada
	if limit <= 0 {
		limit = 10 // Nilai default untuk limit adalah 10
	}

	// Lakukan pengambilan data dari database dengan menggunakan parameter limit, sort, dan sort_by
	var products []models.Product
	query := initialize.DB.Model(&models.Product{})

	// Jika parameter search tidak kosong, tambahkan filter pencarian
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Jika parameter sort dan sort_by disediakan, lakukan pengurutan berdasarkan kolom yang dimaksud
	if sort != "" && sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sort))
	}

	// Limit jumlah data yang diambil sesuai dengan nilai parameter limit
	query = query.Limit(limit)

	// Lakukan pengambilan data
	if err := query.Find(&products).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil produk, kirim respons kesalahan ke klien
		response := helpers.GeneralResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   nil,
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kembalikan respons datatable
	response := helpers.GeneralResponse{
		Code:   200,
		Status: "OK",
		Data:   products,
	}
	return c.JSON(response)
}

func GetProductById(c *fiber.Ctx) error {
	// Ambil ID produk dari parameter URL
	productId, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		// Jika ID tidak valid, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    400,
			Status:  "Bad Request",
			Message: "Id Product tidak valid",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var product models.Product
	if err := initialize.DB.Preload("Category").Preload("File").First(&product, productId).Error; err != nil {
		// Jika produk tidak ditemukan, kirim respons not found ke klien
		response := helpers.ResponseMassage{
			Code:    404,
			Status:  "Not Found",
			Message: "Product Tidak ditemukan",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}
	customResponse := models.ProductResponse{
		ProductId:   product.ProductId,
		Title:       product.Title,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		FileId:      product.FileId,
		CategoryId:  product.CategoryId,
		PathFile:    product.File.Path,         // Mengambil path file dari relasi File
		Category:    product.Category.Category, // Mengambil nama kategori dari relasi Category
	}

	// Kembalikan respons sukses dengan data produk ke klien
	response := helpers.GeneralResponse{
		Code:   200,
		Status: "OK",
		Data:   customResponse,
	}
	return c.JSON(response)
}

func CreateProduct(c *fiber.Ctx) error {
	var requestBody models.ProductCreateRequest
	if err := c.BodyParser(&requestBody); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := validate.Struct(&requestBody); err != nil {
		errors := make(map[string][]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			var tagName string
			switch field {
			case "Title":
				tagName = "title"
			case "Description":
				tagName = "description"
			case "CategoryId":
				tagName = "category_id"
			case "File":
				tagName = "file"
			} // Mendapatkan nama tag JSON yang sesuai
			message := fmt.Sprintf("%s is required", tagName) // Pesan kesalahan yang disesuaikan
			errors[tagName] = append(errors[field], message)
		}
		response := helpers.ResponseError{
			Code:   fiber.StatusBadRequest,
			Status: "Bad Request",
			Error:  errors,
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Ambil data produk dari form
	title := requestBody.Title
	description := requestBody.Description
	categoryId := requestBody.CategoryId

	// Simpan file yang diunggah ke folder public
	file, err := c.FormFile("file")
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "File mohon diisi",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Generate nama unik untuk file yang diunggah
	filename := uuid.New().String() + filepath.Ext(file.Filename)

	// Simpan file ke direktori publik
	if err := c.SaveFile(file, fmt.Sprintf("./public/%s", filename)); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save file",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat entitas File untuk disimpan dalam database
	fileModel := models.File{
		Path:      fmt.Sprintf("/public/%s", filename),
		File_name: filename,
		Size:      strconv.FormatInt(file.Size, 10),
		Format:    filepath.Ext(file.Filename),
	}

	// Simpan file ke dalam database
	if err := initialize.DB.Create(&fileModel).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save file data",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat entitas Product
	product := models.Product{
		Title:       title,
		Description: description,
		CategoryId:  categoryId,
		FileId:      fileModel.FileId,
	}

	// Simpan produk ke dalam database
	if err := initialize.DB.Create(&product).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save product",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kirim respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Product saved successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func EditProduct(c *fiber.Ctx) error {
	var requestBody models.ProductEditRequest
	if err := c.BodyParser(&requestBody); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi struktur data
	if err := validate.Struct(&requestBody); err != nil {
		errors := make(map[string][]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			var tagName string
			switch field {
			case "Title":
				tagName = "title"
			case "ProductId":
				tagName = "product_id"
			case "Description":
				tagName = "description"
			case "CategoryId":
				tagName = "category_id"
			} // Mendapatkan nama tag JSON yang sesuai
			message := fmt.Sprintf("%s is required", tagName) // Pesan kesalahan yang disesuaikan
			errors[tagName] = append(errors[field], message)
		}
		response := helpers.ResponseError{
			Code:   fiber.StatusBadRequest,
			Status: "Bad Request",
			Error:  errors,
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Ambil ID produk dari URL parameter
	productId := requestBody.ProductId

	// Ambil data produk yang akan diedit dari database
	var existingProduct models.Product
	if err := initialize.DB.First(&existingProduct, productId).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusNotFound,
			Status:  "Not Found",
			Message: "Product not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Simpan file yang diunggah ke folder public
	file, err := c.FormFile("file")
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Failed to upload file",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Generate nama unik untuk file yang diunggah
	filename := uuid.New().String() + filepath.Ext(file.Filename)

	// Simpan file ke direktori publik
	if err := c.SaveFile(file, fmt.Sprintf("./public/%s", filename)); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save file",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat entitas File baru untuk disimpan dalam database
	newFile := models.File{
		Path:      fmt.Sprintf("/public/%s", filename),
		File_name: filename,
		Size:      strconv.FormatInt(file.Size, 10),
		Format:    filepath.Ext(file.Filename),
	}

	// Simpan file ke dalam database
	if err := initialize.DB.Create(&newFile).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save new file data",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	if existingProduct.FileId != 0 {
		// Memuat data file terkait
		if err := initialize.DB.Model(&existingProduct).Association("File").Find(&existingProduct.File); err != nil {
			// Handle error jika gagal memuat data file
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Gagal Memuat Data File",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		// Hapus file lama dari sistem file lokal
		oldFile := existingProduct.File
		if err := os.Remove("./" + oldFile.Path); err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: oldFile.Path,
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		// Hapus entitas file lama dari basis data
		if err := initialize.DB.Delete(&oldFile).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Gagal Menghapus Data File",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Update data produk dengan data baru
	existingProduct.Title = requestBody.Title
	existingProduct.Description = requestBody.Description
	existingProduct.CategoryId = requestBody.CategoryId
	existingProduct.FileId = newFile.FileId

	// Simpan data produk yang telah diperbarui ke dalam database
	if err := initialize.DB.Save(&existingProduct).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save updated product data",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kirim respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Product updated successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteProduct(c *fiber.Ctx) error {
	// Ambil ID produk dari parameter URL
	productId, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid product ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Ambil data produk yang akan dihapus dari database
	var product models.Product
	if err := initialize.DB.Preload("File").First(&product, productId).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusNotFound,
			Status:  "Not Found",
			Message: "Product not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Hapus produk dari database
	if err := initialize.DB.Delete(&product).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to delete product",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Jika produk memiliki file terkait, hapus file tersebut
	if err := os.Remove("." + product.File.Path); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Gagal Menghapus Data Dilocal",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Hapus entitas file terkait dari basis data
	if err := initialize.DB.Delete(&product.File).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	// Kirim respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Product Berhasil dihapus",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
