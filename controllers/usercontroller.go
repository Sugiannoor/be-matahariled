package controllers

import (
	"Matahariled/helpers"
	"Matahariled/initialize"
	"Matahariled/models"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
func LoginUser(c *fiber.Ctx) error {
	// Parse body permintaan
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		})
	}

	// Ambil data pengguna dari database berdasarkan email
	var user models.User
	if err := initialize.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
			Code:    fiber.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "Invalid email or password",
		})
	}

	// Periksa kecocokan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
			Code:    fiber.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "Invalid email or password",
		})
	}
	token := jwt.New(jwt.SigningMethodHS256)

	// Set klaim JWT
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.UserId
	claims["username"] = user.UserName
	claims["email"] = user.Email
	claims["phone_number"] = user.PhoneNumber
	claims["fullname"] = user.FullName
	claims["address"] = user.Address
	claims["path"] = user.File.Path
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token berlaku selama 24 jam

	// Tanda tangani token dengan secret key
	secret := []byte(os.Getenv("JWT_SECRET"))
	if secret == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "JWT secret key not found",
		})
	}
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to sign JWT token",
		})
	}
	c.Set("Authorization", "Bearer "+tokenString)
	// Kirim token JWT dan model User dalam respons
	res := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"access_token": tokenString,
			"user":         user,
		},
	}
	return c.JSON(res)
}

func GetUsersLabel(c *fiber.Ctx) error {
	// Ambil semua pengguna dari database
	var users []models.User
	if err := initialize.DB.Find(&users).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil pengguna, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat respons dengan format yang diinginkan
	var userOptions []map[string]interface{}
	for _, user := range users {
		option := map[string]interface{}{
			"value": user.UserId,
			"label": user.FullName, // Atur atribut yang sesuai dengan nama pengguna
		}
		userOptions = append(userOptions, option)
	}

	// Kembalikan respons sukses dengan data pengguna ke klien
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   userOptions,
	}
	return c.JSON(response)
}
func GetCountUser(c *fiber.Ctx) error {
	// Hitung jumlah total pengguna dari database
	var count int64
	if err := initialize.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		// Jika terjadi kesalahan saat menghitung pengguna, kirim respons kesalahan ke klien
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Terjadi Kesalahan Server",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kembalikan respons dengan total pengguna
	response := helpers.GeneralResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   count,
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
func GetProfile(c *fiber.Ctx) error {
	// Ambil ID pengguna dari lokal konteks
	userID, ok := c.Locals("userID").(float64)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  "Internal Server Error",
			"message": "Gagal Mengambil Id Dari Token",
		})
	}

	// Query database untuk mendapatkan profil pengguna berdasarkan ID
	var user models.User
	if err := initialize.DB.Where("user_id = ?", int64(userID)).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"status":  "Not Found",
			"message": "Pengguna Tidak Tersedia",
		})
	}

	// Kirim respons dengan profil pengguna
	return c.JSON(fiber.Map{
		"code":   fiber.StatusOK,
		"status": "OK",
		"data":   user,
	})
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

	validationErrors := make(map[string][]string)

	// Validasi manual untuk setiap field
	if newUser.FullName == "" {
		validationErrors["full_name"] = append(validationErrors["full_name"], "Full Name Wajib Diisi")
	}
	if newUser.UserName == "" {
		validationErrors["username"] = append(validationErrors["username"], "Username Wajib Diisi")
	}
	if newUser.PhoneNumber == "" {
		validationErrors["phone_number"] = append(validationErrors["phone_number"], "Phone Number Wajib Diisi")
	}
	if newUser.Password == "" {
		validationErrors["password"] = append(validationErrors["password"], "Password Wajib Diisi")
	}
	if newUser.Email == "" {
		validationErrors["email"] = append(validationErrors["email"], "Email Wajib Diisi")
	}
	if newUser.Role == "" {
		validationErrors["role"] = append(validationErrors["role"], "Role Wajib Diisi")
	}

	// Jika ada error, kirim response
	if len(validationErrors) > 0 {
		response := helpers.ResponseError{
			Code:   fiber.StatusBadRequest,
			Status: "Bad Request",
			Error:  validationErrors,
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

	user := models.User{
		FullName:    newUser.FullName,
		UserName:    newUser.UserName,
		PhoneNumber: newUser.PhoneNumber,
		Password:    string(hashedPassword),
		Email:       newUser.Email,
		Role:        newUser.Role,
		FileId:      1,
	}
	// Buat pengguna baru di database
	if err := initialize.DB.Create(&user).Error; err != nil {
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
	existingUser.Role = updatedUser.Role
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
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	search := c.Query("search")

	offset := (page - 1) * limit

	// Tentukan default nilai jika parameter tidak ada
	if limit <= 0 {
		limit = 10 // Misalnya, nilai default untuk limit adalah 10
	}

	// Lakukan pengambilan data dari database dengan menggunakan parameter limit, sort, dan sort_by
	var users []models.UserResponse
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
	query = query.Limit(limit).Offset(offset)

	var totalRecords int64
	if err := initialize.DB.Model(&models.Product{}).Count(&totalRecords).Error; err != nil {
		response := helpers.GeneralResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   nil,
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	// Limit jumlah data yang diambil sesuai dengan nilai parameter limit

	if err := query.Find(&users).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil data, kirim respons kesalahan ke klien
		response := helpers.GeneralResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   []map[string]interface{}{},
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	// Kembalikan respons datatable
	response := helpers.DataTableResponse{
		CurrentPage:  1,                               // Nomor halaman saat ini (default 1)
		FirstPageURL: "",                              // URL halaman pertama
		From:         1,                               // Nomor record pertama pada halaman saat ini
		LastPage:     totalPages,                      // Total jumlah halaman (default 1)
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

func CreateUserForm(c *fiber.Ctx) error {
	// Ambil data pengguna dari form
	fullName := c.FormValue("full_name")
	userName := c.FormValue("username")
	phoneNumber := c.FormValue("phone_number")
	password := c.FormValue("password")
	email := c.FormValue("email")
	address := c.FormValue("address")
	role := c.FormValue("role")

	var fileModel models.File

	// Periksa apakah ada file yang diunggah
	file, err := c.FormFile("file")
	if err == nil {
		// Jika ada file yang diunggah
		filename := uuid.New().String() + filepath.Ext(file.Filename)

		// Buat entitas File untuk disimpan dalam database
		fileModel = models.File{
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
				Message: "Failed to save file database",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
		if err := c.SaveFile(file, fmt.Sprintf("./public/%s", filename)); err != nil {
			response := helpers.ResponseMassage{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to save file",
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to hash password",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Buat entitas pengguna
	user := models.User{
		FullName:    fullName,
		UserName:    userName,
		PhoneNumber: phoneNumber,
		Password:    string(hashedPassword),
		Email:       email,
		Address:     &address,
		Role:        role,
	}

	// Tambahkan ID file jika ada file yang diunggah
	if fileModel.FileId != 0 {
		user.FileId = fileModel.FileId
	}
	// Simpan pengguna ke dalam database
	if err := initialize.DB.Create(&user).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to save user",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Kirim respons sukses
	response := helpers.ResponseMassage{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "User saved successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteUserT(c *fiber.Ctx) error {
	UserId, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid User ID",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Ambil data produk yang akan dihapus dari database
	var User models.User
	if err := initialize.DB.Preload("File").First(&User, UserId).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusNotFound,
			Status:  "Not Found",
			Message: "User not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Hapus produk dari database
	if err := initialize.DB.Delete(&User).Error; err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to delete User",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Jika produk memiliki file terkait, hapus file tersebut
	if err := os.Remove("." + User.File.Path); err != nil {
		response := helpers.ResponseMassage{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Gagal Menghapus Data Dilocal",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Hapus entitas file terkait dari basis data
	if err := initialize.DB.Delete(&User.File).Error; err != nil {
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
		Message: "User Berhasil dihapus",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
