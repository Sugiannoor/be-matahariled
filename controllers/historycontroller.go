package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllHistories(c *fiber.Ctx) error {
	// Ambil semua data history dari database
	var histories []models.History
	if err := initialize.DB.Preload("Product").Preload("Product.Category").Preload("File").Find(&histories).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil history, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to fetch histories",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Membuat slice untuk menyimpan respons history
	historyResponses := make([]models.HistoryResponse, len(histories))

	// Mengisi data respons history
	for i, history := range histories {
		historyResponses[i] = models.HistoryResponse{
			HistoryId:    history.HistoryId,
			Title:        history.Title,
			Description:  history.Description,
			ProductName:  history.Product.Title,
			CategoryName: history.Product.Category.Category,
			PathFile:     history.File.Path,
			CreatedAt:    history.CreatedAt,
			UpdatedAt:    history.UpdatedAt,
		}
	}

	// Mengirimkan respons sukses dengan daftar history
	response := helpers.ResponseGetAll{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   historyResponses,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func GetDatatableHistories(c *fiber.Ctx) error {
	// Ambil nilai parameter limit, sort, sort_by, search, dan product_id dari query string
	limit, _ := strconv.Atoi(c.Query("limit"))
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	search := c.Query("search")
	productID := c.Query("product_id")

	// Tentukan default nilai jika parameter tidak ada
	if limit <= 0 {
		limit = 10 // Nilai default untuk limit adalah 10
	}

	// Lakukan pengambilan data dari database dengan menggunakan parameter limit, sort, dan sort_by
	var histories []models.History
	query := initialize.DB.Preload("Product").Preload("Product.Category").Preload("File").Model(&models.History{})

	// Jika parameter search tidak kosong, tambahkan filter pencarian
	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Jika parameter product_id disediakan, lakukan filter berdasarkan product_id
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	// Jika parameter sort dan sort_by disediakan, lakukan pengurutan berdasarkan kolom yang dimaksud
	if sort != "" && sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sort))
	}

	// Limit jumlah data yang diambil sesuai dengan nilai parameter limit
	query = query.Limit(limit)

	// Lakukan pengambilan data
	if err := query.Find(&histories).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil history, kirim respons kesalahan ke klien
		response := helpers.GeneralResponse{
			Code:   fiber.StatusInternalServerError,
			Status: "Internal Server Error",
			Data:   nil,
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Siapkan respons DataTable
	response := helpers.DataTableResponse{
		CurrentPage:  1,                                   // Nomor halaman saat ini (default 1)
		FirstPageURL: "",                                  // URL halaman pertama
		From:         1,                                   // Nomor record pertama pada halaman saat ini
		LastPage:     1,                                   // Total jumlah halaman (default 1)
		LastPageURL:  "",                                  // URL halaman terakhir
		NextPageURL:  "",                                  // URL halaman berikutnya
		PrevPageURL:  "",                                  // URL halaman sebelumnya
		To:           len(histories),                      // Nomor record terakhir pada halaman saat ini
		Total:        len(histories),                      // Total jumlah record
		Data:         make([]interface{}, len(histories)), // Buat slice interface{} dengan panjang yang sama dengan histories
	}

	// Mengisi data respons DataTable
	for i, history := range histories {
		// Buat map untuk setiap history
		historyMap := map[string]interface{}{
			"history_id":    history.HistoryId,
			"title":         history.Title,
			"description":   history.Description,
			"product_id":    history.ProductId,
			"product_name":  history.Product.Title,
			"category_name": history.Product.Category.Category,
			"file_id":       history.FileId,
			"path_file":     history.File.Path,
			"created_at":    history.CreatedAt,
			"updated_at":    history.UpdatedAt,
		}

		// Tambahkan map history ke dalam slice Data pada respons
		response.Data[i] = historyMap
	}

	// Kembalikan respons JSON dengan format DataTable
	return c.JSON(helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func CreateHistory(c *fiber.Ctx) error {
	// Ambil data produk dari form
	title := c.FormValue("title")
	description := c.FormValue("description")
	productIdStr := c.FormValue("product_id")
	productId, err := strconv.ParseInt(productIdStr, 10, 64)
	if err != nil || productId <= 0 {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid Product Id",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var product models.Product
	if err := initialize.DB.First(&product, productId).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Product Tidak Tersedia",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if title == "" || description == "" {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Judul atau Deskripsi diperlukan",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

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
	history := models.History{
		Title:       title,
		Description: description,
		ProductId:   productId,
		FileId:      fileModel.FileId,
	}

	// Simpan produk ke dalam database
	if err := initialize.DB.Create(&history).Error; err != nil {
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
		Message: "History saved successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func UpdateHistory(c *fiber.Ctx) error {
	// Ambil history ID dari parameter URL
	historyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || historyID <= 0 {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid or missing History ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Periksa apakah riwayat dengan historyID tersebut ada di database
	var history models.History
	if err := initialize.DB.First(&history, historyID).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "History Not Found",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Ambil data baru dari form
	title := c.FormValue("title")
	description := c.FormValue("description")
	productIdStr := c.FormValue("product_id")
	productId, err := strconv.ParseInt(productIdStr, 10, 64)
	if err != nil || productId <= 0 {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid or missing Product Id",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi apakah title atau description kosong
	if title == "" || description == "" {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Title and description are required",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Periksa apakah product dengan productId tersebut ada di database
	var product models.Product
	if err := initialize.DB.First(&product, productId).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Product with given ID not found",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

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
		if history.FileId != 0 {
			// Memuat data file terkait
			if err := initialize.DB.Model(&history).Association("File").Find(&history.File); err != nil {
				// Handle error jika gagal memuat data file
				response := helpers.ResponseMassage{
					Code:    fiber.StatusInternalServerError,
					Status:  "Internal Server Error",
					Message: "Gagal Memuat Data File",
				}
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			// Hapus file lama dari sistem file lokal
			oldFile := history.File
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
		history.File = newFile
	}

	// Update data riwayat dengan data baru
	history.Title = title
	history.Description = description
	history.ProductId = productId

	// Simpan perubahan ke dalam database
	if err := initialize.DB.Save(&history).Error; err != nil {
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
func DeleteHistory(c *fiber.Ctx) error {
	// Ambil history ID dari parameter URL
	historyID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || historyID <= 0 {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid or missing History ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Periksa apakah riwayat dengan historyID tersebut ada di database
	var history models.History
	if err := initialize.DB.First(&history, historyID).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "History Not Found",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Hapus file terkait dari sistem file lokal
	if history.FileId != 0 {
		// Memuat data file terkait
		if err := initialize.DB.Model(&history).Association("File").Find(&history.File); err != nil {
			// Handle error jika gagal memuat data file
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to load file data",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		// Hapus file dari sistem file lokal
		if err := os.Remove("." + history.File.Path); err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to delete local file",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		// Hapus entitas file dari basis data
		if err := initialize.DB.Delete(&history.File).Error; err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to delete file data",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Hapus riwayat dari basis data
	if err := initialize.DB.Delete(&history).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to delete history",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kirim respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "History deleted successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
