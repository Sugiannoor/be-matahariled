package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllVideos(c *fiber.Ctx) error {
	// Ambil semua video dari database
	var videos []models.Video
	if err := initialize.DB.Find(&videos).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to fetch videos",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   videos,
	}
	return c.JSON(response)
}

func GetDatatableVideos(c *fiber.Ctx) error {
	// Ambil nilai parameter limit, sort, sort_by, dan search dari query string
	limit, _ := strconv.Atoi(c.Query("limit"))
	sort := c.Query("sort")
	page, _ := strconv.Atoi(c.Query("page"))
	sortBy := c.Query("sort_by")
	search := c.Query("search")

	offset := (page - 1) * limit
	// Tentukan default nilai jika parameter tidak ada
	if limit <= 0 {
		limit = 10 // Nilai default untuk limit adalah 10
	}

	// Lakukan pengambilan data dari database dengan menggunakan parameter limit, sort, dan sort_by
	var videos []models.Video
	query := initialize.DB.Model(&models.Video{})

	// Jika parameter search tidak kosong, tambahkan filter pencarian
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Jika parameter sort dan sort_by disediakan, lakukan pengurutan berdasarkan kolom yang dimaksud
	if sort != "" && sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sort))
	}

	var totalRecords int64
	if err := initialize.DB.Model(&models.Product{}).Count(&totalRecords).Error; err != nil {
		response := helpers.GeneralResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   []map[string]interface{}{},
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	// Limit jumlah data yang diambil sesuai dengan nilai parameter limit
	query = query.Limit(limit).Offset(offset)

	// Lakukan pengambilan data
	if err := query.Find(&videos).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil data video, kirim respons kesalahan ke klien
		response := helpers.GeneralResponse{
			Code:   fiber.StatusInternalServerError,
			Status: "Internal Server Error",
			Data:   []map[string]interface{}{},
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	// Siapkan respons DataTable
	response := helpers.DataTableResponse{
		CurrentPage:  1,                                // Nomor halaman saat ini (default 1)
		FirstPageURL: "",                               // URL halaman pertama
		From:         1,                                // Nomor record pertama pada halaman saat ini
		LastPage:     totalPages,                       // Total jumlah halaman (default 1)
		LastPageURL:  "",                               // URL halaman terakhir
		NextPageURL:  "",                               // URL halaman berikutnya
		PrevPageURL:  "",                               // URL halaman sebelumnya
		To:           len(videos),                      // Nomor record terakhir pada halaman saat ini
		Total:        len(videos),                      // Total jumlah record
		Data:         make([]interface{}, len(videos)), // Buat slice interface{} dengan panjang yang sama dengan videos
	}

	// Mengisi data respons DataTable
	for i, video := range videos {
		// Buat map untuk setiap video
		videoMap := map[string]interface{}{
			"video_id":    video.VideoId,
			"video_title": video.Title,
			"embed":       video.Embed,
			"created_at":  video.CreatedAt,
			"updated_at":  video.UpdatedAt,
		}

		// Tambahkan map video ke dalam slice Data pada respons
		response.Data[i] = videoMap
	}

	// Kembalikan respons JSON dengan format DataTable
	return c.JSON(helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   response,
	})
}
