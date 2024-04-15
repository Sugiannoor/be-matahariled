package middleware

import (
	"Matahariled/helpers"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware adalah middleware untuk melakukan autentikasi JWT
func JWTMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Missing Authorization header",
			})
		}

		// Parse token
		token := strings.Split(header, " ")
		if len(token) != 2 || token[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Invalid authorization token",
			})
		}

		// Validasi token
		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token[1], &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
					Code:    fiber.StatusUnauthorized,
					Status:  "Unauthorized",
					Message: "Invalid token signature",
				})
			}
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseMassage{
				Code:    fiber.StatusBadRequest,
				Status:  "Bad Request",
				Message: "Failed to parse JWT token",
			})
		}
		if !parsedToken.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseMassage{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Invalid token",
			})
		}

		// Set data pengguna dari token ke konteks jika valid
		c.Locals("user_id", claims["user_id"])
		c.Locals("email", claims["email"])
		c.Locals("role", claims["role"])

		// Periksa peran pengguna jika peran yang diperlukan diberikan
		if len(allowedRoles) > 0 {
			role := claims["role"].(string)
			roleAllowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				return c.Status(fiber.StatusForbidden).JSON(helpers.ResponseMassage{
					Code:    fiber.StatusForbidden,
					Status:  "Forbidden",
					Message: "User role is not allowed to access this resource",
				})
			}
		}

		return c.Next()
	}
}
