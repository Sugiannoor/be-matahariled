package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
			ProductId:     product.ProductId,
			Title:         product.Title,
			Description:   product.Description,
			Specification: product.Specification,
			CreatedAt:     product.CreatedAt,
			UpdatedAt:     product.UpdatedAt,
			FileId:        product.FileId,
			CategoryId:    product.CategoryId,
			PathFile:      product.File.Path,
			Category:      product.Category.Category,
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

func GetProductsLabel(c *fiber.Ctx) error {
	// Ambil semua produk dari database
	var products []models.Product
	if err := initialize.DB.Find(&products).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil produk, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat respons dengan format yang diinginkan
	var productOptions []map[string]interface{}
	for _, product := range products {
		option := map[string]interface{}{
			"value": product.ProductId,
			"label": product.Title,
		}
		productOptions = append(productOptions, option)
	}

	// Kembalikan respons sukses dengan data produk ke klien
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   productOptions,
	}
	return c.JSON(response)
}
func GetCountProduct(c *fiber.Ctx) error {
	// Hitung jumlah total produk dari database
	var count int64
	if err := initialize.DB.Model(&models.Product{}).Count(&count).Error; err != nil {
		// Jika terjadi kesalahan saat menghitung produk, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kembalikan respons dengan total produk
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   count,
	}
	return c.JSON(response)
}

func GetDatatableProducts(c *fiber.Ctx) error {
	// Ambil nilai parameter limit, page, sort, sort_by, dan search dari query string
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	search := c.Query("search")
	categoryID := c.Query("category_id")

	// Tentukan default nilai jika parameter tidak ada
	if limit <= 0 {
		limit = 10 // Nilai default untuk limit adalah 10
	}
	if page <= 0 {
		page = 1 // Halaman default adalah 1
	}

	// Hitung offset berdasarkan halaman dan limit
	offset := (page - 1) * limit

	// Lakukan pengambilan data dari database dengan menggunakan parameter limit, offset, sort, dan sort_by
	var products []models.Product
	query := initialize.DB.Preload("File").Preload("Category").Model(&models.Product{})

	// Jika parameter search tidak kosong, tambahkan filter pencarian
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Jika parameter sort dan sort_by disediakan, lakukan pengurutan berdasarkan kolom yang dimaksud
	if sort != "" && sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sort))
	}

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// Limit jumlah data yang diambil sesuai dengan nilai parameter limit dan offset
	query = query.Limit(limit).Offset(offset)

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

	// Hitung total jumlah record tanpa paginasi
	var totalRecords int64
	if err := initialize.DB.Model(&models.Product{}).Count(&totalRecords).Error; err != nil {
		response := helpers.GeneralResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   nil,
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Hitung total jumlah halaman berdasarkan total jumlah record dan limit
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	response := helpers.DataTableResponse{
		CurrentPage:  page,
		FirstPageURL: "", // Anda bisa menentukan URL halaman pertama jika perlu
		From:         offset + 1,
		LastPage:     totalPages,
		LastPageURL:  "", // Anda bisa menentukan URL halaman terakhir jika perlu
		NextPageURL:  "", // Anda bisa menentukan URL halaman berikutnya jika perlu
		PrevPageURL:  "", // Anda bisa menentukan URL halaman sebelumnya jika perlu
		To:           offset + len(products),
		Total:        int(totalRecords),
		Data:         make([]interface{}, len(products)),
	}

	for i, product := range products {
		// Buat map untuk setiap produk
		productMap := map[string]interface{}{
			"product_id":    product.ProductId,
			"name":          product.Title,
			"description":   product.Description,
			"specification": product.Specification,
			"created_at":    product.CreatedAt,
			"updated_at":    product.UpdatedAt,
			"file_id":       product.FileId,
			"category_id":   product.CategoryId,
			"path_file":     product.File.Path,
			"category":      product.Category.Category,
		}

		// Tambahkan map produk ke dalam slice Data pada respons
		response.Data[i] = productMap
	}

	// Kembalikan respons JSON dengan format datatable
	return c.JSON(helpers.GeneralResponse{
		Code:   200,
		Status: "OK",
		Data:   response,
	})
}

func GetProductById(c *fiber.Ctx) error {
	// Ambil ID produk dari parameter URL
	productId := c.Params("id")

	// Buat variabel untuk menyimpan data produk
	var product models.Product

	// Cari produk berdasarkan ID
	if err := initialize.DB.Preload("Category").Preload("File").Where("product_id = ?", productId).First(&product).Error; err != nil {
		// Jika produk tidak ditemukan, kirim respons not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusNotFound,
				Status:  "Not Found",
				Message: "Product not found",
			}
			return c.Status(fiber.StatusNotFound).JSON(response)
		}
		// Jika terjadi kesalahan lain saat mengambil produk, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to fetch product",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Membuat respons untuk produk yang ditemukan
	productResponse := models.ProductResponse{
		ProductId:     product.ProductId,
		Title:         product.Title,
		Specification: product.Specification,
		Description:   product.Description,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
		FileId:        product.FileId,
		CategoryId:    product.CategoryId,
		PathFile:      product.File.Path,
		Category:      product.Category.Category,
	}

	// Mengirimkan respons sukses dengan data produk yang ditemukan
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   productResponse,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateProduct(c *fiber.Ctx) error {
	// Ambil data produk dari form
	title := c.FormValue("name")
	description := c.FormValue("description")
	specification := c.FormValue("specification")
	categoryId, err := strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid category ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Simpan file yang diunggah ke folder public
	file, err := c.FormFile("file")
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "File is required",
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
		Title:         title,
		Specification: specification,
		Description:   description,
		CategoryId:    categoryId,
		FileId:        fileModel.FileId,
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

func UpdateProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || productID <= 0 {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid or missing product ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Periksa apakah riwayat dengan productID tersebut ada di database
	var product models.Product
	if err := initialize.DB.First(&product, productID).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Product Not Found",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Ambil data baru dari form
	title := c.FormValue("title")
	description := c.FormValue("description")
	specification := c.FormValue("specification")
	categoryIdStr := c.FormValue("category_id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil || categoryId <= 0 {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid or missing Product Id",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi apakah title atau description kosong
	if title == "" || description == "" || specification == "" {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Title and description are required",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Periksa apakah product dengan productId tersebut ada di database

	// Cek apakah ada file baru yang diunggah
	file, err := c.FormFile("file")
	if err != nil {
	} else {
		// Jika ada file baru, simpan file baru dan hapus file lama
		// Generate nama unik untuk file yang diunggah
		filename := uuid.New().String() + filepath.Ext(file.Filename)
		// eapatkan data file id  dari baris tesb
		// Simpan file ke direktori publik
		if err := c.SaveFile(file, fmt.Sprintf("./public/%s", filename)); err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to save file",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		// Hapus file lama dari direktori publik
		if product.FileId != 0 {
			// Memuat data file terkait
			if err := initialize.DB.Model(&product).Association("File").Find(&product.File); err != nil {
				// Handle error jika gagal memuat data file
				response := helpers.ResponseMassage{
					Code:    fiber.StatusInternalServerError,
					Status:  "Internal Server Error",
					Message: "Gagal Memuat Data File",
				}
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			// Hapus file lama dari sistem file lokal
			oldFile := product.File
			if err := os.Remove("./" + oldFile.Path); err != nil {
				response := helpers.ResponseMassage{
					Code:    fiber.StatusInternalServerError,
					Status:  "Internal Server Error",
					Message: "Gagal Menghapus Data Local",
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

		// Buat entitas File baru untuk disimpan dalam database
		newFile := models.File{
			Path:      fmt.Sprintf("/public/%s", filename),
			File_name: filename,
			Size:      strconv.FormatInt(file.Size, 10),
			Format:    filepath.Ext(file.Filename),
		}

		// Simpan file baru ke dalam database
		if err := initialize.DB.Create(&newFile).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to save new file data",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		// Ganti file lama dengan file baru
		product.File = newFile
	}

	// Update data riwayat dengan data baru
	product.Title = title
	product.Specification = specification
	product.Description = description
	product.ProductId = productID
	product.CategoryId = categoryId

	// Simpan perubahan ke dalam database
	if err := initialize.DB.Save(&product).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to update history",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kirim respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "History updated successfully",
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

func CreateProductT(c *fiber.Ctx) error {
	// Ambil data produk dari form
	title := c.FormValue("name")
	description := c.FormValue("description")
	specification := c.FormValue("specification")
	categoryId, err := strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid category ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Simpan file utama yang diunggah ke folder public
	file, err := c.FormFile("file")
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Main file is required",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Generate nama unik untuk file utama yang diunggah
	mainFilename := uuid.New().String() + filepath.Ext(file.Filename)

	// Simpan file utama ke direktori publik
	if err := c.SaveFile(file, fmt.Sprintf("./public/%s", mainFilename)); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save main file",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat entitas File untuk file utama yang akan disimpan dalam database
	mainFileModel := models.File{
		Path:      fmt.Sprintf("/public/%s", mainFilename),
		File_name: mainFilename,
		Size:      strconv.FormatInt(file.Size, 10),
		Format:    filepath.Ext(file.Filename),
	}

	// Simpan file utama ke dalam database
	if err := initialize.DB.Create(&mainFileModel).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save main file data",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat entitas Product
	product := models.Product{
		Title:         title,
		Specification: specification,
		Description:   description,
		CategoryId:    categoryId,
		FileId:        mainFileModel.FileId,
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

	// Simpan galeri ke dalam database
	form, err := c.MultipartForm()
	if err != nil {
	} else {

		files := form.File["gallery[]"]
		var galleries []models.Gallery
		for _, file := range files {
			// Generate nama unik untuk file yang diunggah
			filename := uuid.New().String() + filepath.Ext(file.Filename)

			// Simpan file ke direktori publik
			if err := c.SaveFile(file, fmt.Sprintf("./public/%s", filename)); err != nil {
				response := helpers.ResponseMassage{
					Code:    fiber.StatusInternalServerError,
					Status:  "Internal Server Error",
					Message: "Failed to save gallery files",
				}
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			// Buat entitas Gallery untuk disimpan dalam database
			gallery := models.Gallery{
				Path:         fmt.Sprintf("/public/%s", filename),
				Gallery_name: file.Filename,
				Size:         strconv.FormatInt(file.Size, 10),
				Format:       filepath.Ext(file.Filename),
				ProductId:    product.ProductId, // Gunakan ID produk yang baru dibuat
			}

			galleries = append(galleries, gallery)
		}

		// Simpan galeri ke dalam database
		if err := initialize.DB.Create(&galleries).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to save gallery data",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Kirim respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Product and gallery files saved successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteProductT(c *fiber.Ctx) error {
	// Ambil ID produk dari parameter URL
	var hero models.Hero
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
	if err := initialize.DB.Preload("Gallery").Preload("File").First(&product, productId).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusNotFound,
			Status:  "Not Found",
			Message: "Product not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Hapus galeri-galeri terkait dengan produk
	for _, gallery := range product.Gallery {
		if err := os.Remove("." + gallery.Path); err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to delete gallery",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
		if err := initialize.DB.Delete(&gallery).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to delete gallery",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	if err := os.Remove("." + product.File.Path); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to delete file",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	if err := initialize.DB.Delete(&product.File).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to delete file",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	if hero.ProductId != 0 {
		if err := os.Remove("." + hero.Path); err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to delete gallery",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
		if err := initialize.DB.Delete(&hero).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to delete product",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
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

	// Kirim respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Product and its galleries deleted successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
