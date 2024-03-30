package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Index(c *fiber.Ctx) error {
	var user []models.User
	initialize.DB.Find(&user)

	response := helpers.ResponseGetAll{
		Code:   200,
		Status: "OK",
		Data:   user,
	}

	return c.JSON(response)
}
func GetUserById(c *fiber.Ctx) error {
	// Ambil ID pengguna dari parameter URL
	userIdStr := c.Query("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		// Jika ID tidak valid, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    400,
			Status:  "Bad Request",
			Message: "ID pengguna tidak valid",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Buat variabel untuk menyimpan informasi pengguna
	var user models.User

	// Cari pengguna dengan ID yang sesuai dalam database
	if err := initialize.DB.Where("user_id = ?", userId).First(&user).Error; err != nil {
		// Jika pengguna tidak ditemukan, kirim respons ke klien
		response := helpers.ResponseMassage{
			Code:    404,
			Status:  "Not Found",
			Message: "Pengguna tidak ditemukan",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Jika pengguna ditemukan, kirim respons dengan data pengguna ke klien
	response := helpers.ResponseGetSingle{
		Code:   200,
		Status: "OK",
		Data:   user,
	}
	return c.JSON(response)
}

// Create
func CreateUser(c *fiber.Ctx) error {
	// Ambil data yang diterima dari permintaan
	var newUser models.User
	if err := c.BodyParser(&newUser); err != nil {
		// Jika terjadi kesalahan dalam mengurai permintaan, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    400,
			Status:  "Bad Request",
			Message: "Gagal Membuat User",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi data
	if err := validate.Struct(newUser); err != nil {
		// Jika terjadi kesalahan validasi, kirim respons kesalahan ke klien
		// Format respons bad request dengan detail pesan kesalahan
		errors := make(map[string][]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			var tagName string
			switch field {
			case "FullName":
				tagName = "full_name"
			case "UserId":
				tagName = "user_id"
			case "UserName":
				tagName = "username"
			case "PhoneNumber":
				tagName = "phone_number"
			case "Password":
				tagName = "password"
			case "Email":
				tagName = "email"
			case "Address":
				tagName = "address"
			case "Role":
				tagName = "role"
			case "CreatedAt":
				tagName = "created_at"
			case "UpdatedAt":
				tagName = "updated_at"
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

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		// Jika terjadi kesalahan dalam menghasilkan hash password, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Simpan hash password ke field Password
	newUser.Password = string(hashedPassword)

	// Buat pengguna baru di database
	if err := initialize.DB.Create(&newUser).Error; err != nil {
		// Jika terjadi kesalahan dalam membuat pengguna baru, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Jika pembuatan pengguna berhasil, kirim respons sukses ke klien
	response := helpers.ResponseMassage{
		Code:    200,
		Status:  "OK",
		Message: "Pengguna berhasil dibuat",
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// Delete

func DeleteUser(c *fiber.Ctx) error {
	// Ambil ID pengguna dari query string
	userIdStr := c.Query("id")

	// Konversi ID pengguna dari string ke integer
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    400,
			Status:  "Bad Request",
			Message: "ID pengguna tidak valid",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Buat variabel untuk menyimpan informasi pengguna yang akan dihapus
	var user models.User

	// Cari pengguna dengan ID yang sesuai dalam database
	if err := initialize.DB.Where("user_id = ?", userId).First(&user).Error; err != nil {
		// Jika pengguna tidak ditemukan, kirim respons ke klien
		response := helpers.ResponseMassage{
			Code:    404,
			Status:  "Not Found",
			Message: "Pengguna tidak ditemukan",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Hapus pengguna dari database
	if err := initialize.DB.Delete(&user).Error; err != nil {
		// Jika terjadi kesalahan saat menghapus, kirim respons ke klien
		response := helpers.ResponseMassage{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Jika pengguna berhasil dihapus, kirim respons sukses ke klien
	response := helpers.ResponseMassage{
		Code:    200,
		Status:  "OK",
		Message: "Pengguna berhasil dihapus",
	}
	return c.JSON(response)
}

func EditUser(c *fiber.Ctx) error {
	// Ambil data pengguna yang akan diubah dari body permintaan
	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		response := helpers.ResponseMassage{
			Code:    400,
			Status:  "Bad Request",
			Message: "Gagal mengurai data pengguna yang diperbarui",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Cari pengguna dengan ID yang sesuai dalam database
	var existingUser models.User
	if err := initialize.DB.Where("user_id = ?", updatedUser.UserId).First(&existingUser).Error; err != nil {
		// Jika pengguna tidak ditemukan, kirim respons ke klien
		response := helpers.ResponseMassage{
			Code:    404,
			Status:  "Not Found",
			Message: "Pengguna tidak ditemukan",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Update data pengguna yang ada dengan data baru
	existingUser.FullName = updatedUser.FullName
	existingUser.UserName = updatedUser.UserName
	existingUser.Email = updatedUser.Email
	existingUser.PhoneNumber = updatedUser.PhoneNumber
	existingUser.Address = updatedUser.Address
	// Jika diperlukan, Anda dapat menambahkan logika untuk mengedit bidang lainnya

	// Simpan perubahan ke database
	if err := initialize.DB.Save(&existingUser).Error; err != nil {
		// Jika terjadi kesalahan saat menyimpan, kirim respons ke klien
		response := helpers.ResponseMassage{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server saat menyimpan perubahan",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Jika pengguna berhasil diubah, kirim respons sukses ke klien
	response := helpers.ResponseMassage{
		Code:    200,
		Status:  "OK",
		Message: "Pengguna berhasil diperbarui",
	}
	return c.JSON(response)
}

func UserDatatable(c *fiber.Ctx) error {
	// Ambil nilai parameter limit, sort, sort_by, dan search dari query string
	limit, _ := strconv.Atoi(c.Query("limit"))
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	search := c.Query("search")

	// Tentukan default nilai jika parameter tidak ada
	if limit <= 0 {
		limit = 10 // Misalnya, nilai default untuk limit adalah 10
	}

	// Lakukan pengambilan data dari database dengan menggunakan parameter limit, sort, dan sort_by
	var users []models.User
	query := initialize.DB.Model(&models.User{})

	// Jika parameter search tidak kosong, tambahkan filter pencarian
	if search != "" {
		query = query.Where("user_name LIKE ? OR full_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Jika parameter sort dan sort_by disediakan, lakukan pengurutan berdasarkan kolom yang dimaksud
	if sort != "" && sortBy != "" {
		// Jika sort adalah "asc", lakukan pengurutan menaik berdasarkan kolom sort_by
		if sort == "asc" {
			query = query.Order(fmt.Sprintf("%s %s", sortBy, "ASC"))
		} else if sort == "desc" {
			// Jika sort adalah "desc", lakukan pengurutan menurun berdasarkan kolom sort_by
			query = query.Order(fmt.Sprintf("%s %s", sortBy, "DESC"))
		}
	}

	// Limit jumlah data yang diambil sesuai dengan nilai parameter limit
	query = query.Limit(limit)

	// Lakukan pengambilan data
	query.Find(&users)

	// Kembalikan respons datatable
	response := helpers.DataTableResponse{
		CurrentPage:  1,                               // Nomor halaman saat ini (default 1)
		FirstPageURL: "",                              // URL halaman pertama
		From:         1,                               // Nomor record pertama pada halaman saat ini
		LastPage:     1,                               // Total jumlah halaman (default 1)
		LastPageURL:  "",                              // URL halaman terakhir
		NextPageURL:  "",                              // URL halaman berikutnya
		PrevPageURL:  "",                              // URL halaman sebelumnya
		To:           len(users),                      // Nomor record terakhir pada halaman saat ini
		Total:        len(users),                      // Total jumlah record
		Data:         make([]interface{}, len(users)), // Buat slice interface{} dengan panjang yang sama dengan users
	}

	// Masukkan data dari users ke dalam respons
	for i, user := range users {
		response.Data[i] = user
	}

	// Kembalikan respons umum dengan data datatable
	return c.JSON(helpers.GeneralResponse{
		Code:   200,
		Status: "OK",
		Data:   response,
	})
}
