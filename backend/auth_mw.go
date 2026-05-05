package main

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	EmployeeID string
	Role       string
}

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")

		if !strings.HasPrefix(auth, "Bearer ") {
			return fiber.NewError(401, "missing token")
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return fiber.NewError(401, "invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)

		user := &User{
			EmployeeID: claims["employee_id"].(string),
			Role:       claims["role"].(string),
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func ManagerOnly(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		var exists bool
		err := pool.QueryRow(c.Context(),
			`SELECT EXISTS(SELECT 1 FROM managers WHERE employee_id=$1)`,
			user.EmployeeID,
		).Scan(&exists)

		if err != nil || !exists {
			return fiber.NewError(403, "not a manager")
		}

		return c.Next()
	}
}
