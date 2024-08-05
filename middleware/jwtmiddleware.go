package middleware

import (
	"Matahariled/helpers"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func MultiRoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		tokenString := c.Get("Authorization")

		// Periksa apakah token kosong
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
				Code:    fiber.StatusBadRequest,
				Status:  "Bad Request",
				Message: "Token Diperlukan",
			})
		}

		// Pisahkan nilai token JWT dari string "Bearer <token>"
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Invalid Token Bearrer",
			})
		}
		jwtToken := tokenParts[1]

		var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
		// Parse token JWT
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak valid: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Terjadi Kesalahan",
			})
		}

		// Periksa apakah token valid
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Token Tidak Valid",
			})
		}

		// Ambil klaim dari token
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("userID", claims["user_id"])
		// Periksa apakah peran pengguna ada di dalam daftar peran yang diperbolehkan
		userRole := claims["role"].(string)
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				roleAllowed = true
				break
			}
		}

		// Jika peran pengguna tidak diizinkan, kembalikan status Unauthorized
		if !roleAllowed {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Anda Tidak Ada Akses",
			})
		}

		// Jika peran pengguna diizinkan, lanjutkan ke handler berikutnya
		return c.Next()
	}
}
