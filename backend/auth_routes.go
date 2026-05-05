package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	v := os.Getenv("JWT_SECRET")
	if v == "" {
		return "dev-secret-change-this"
	}
	return v
}

type AuthRegisterReq struct {
	EmployeeID string `json:"employee_id"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
	Team       string `json:"team"`
	Password   string `json:"password"`
}

type AuthLoginReq struct {
	EmployeeID string `json:"employee_id"`
	Password   string `json:"password"`
}

type AuthUserResponse struct {
	EmployeeID string `json:"employee_id"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
	Team       string `json:"team"`
	IsManager  bool   `json:"is_manager"`
}

type AuthResponse struct {
	Token string           `json:"token"`
	User  AuthUserResponse `json:"user"`
}

func createToken(employeeID string, role string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"employee_id": employeeID,
		"role":        role,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	})
	return t.SignedString(jwtSecret)
}

func getAuthUser(pool *pgxpool.Pool, employeeID string) (AuthUserResponse, error) {
	var u AuthUserResponse

	err := pool.QueryRow(context.Background(), `
		SELECT
			e.employee_id,
			e.full_name,
			e.role,
			e.team,
			EXISTS(SELECT 1 FROM managers m WHERE m.employee_id = e.employee_id) AS is_manager
		FROM employees e
		WHERE e.employee_id = $1
	`, employeeID).Scan(&u.EmployeeID, &u.FullName, &u.Role, &u.Team, &u.IsManager)

	return u, err
}

func registerAuthRoutes(app *fiber.App, pool *pgxpool.Pool) {
	app.Post("/v1/auth/register", func(c *fiber.Ctx) error {
		var req AuthRegisterReq
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(400, err.Error())
		}

		req.EmployeeID = strings.TrimSpace(req.EmployeeID)
		req.FullName = strings.TrimSpace(req.FullName)
		req.Role = strings.TrimSpace(req.Role)
		req.Team = strings.TrimSpace(req.Team)

		if req.EmployeeID == "" || req.FullName == "" || req.Role == "" || req.Team == "" || req.Password == "" {
			return fiber.NewError(400, "employee_id, full_name, role, team, password are required")
		}
		if req.Role != "ux" && req.Role != "qa" && req.Role != "ba" {
			return fiber.NewError(400, "role must be ux, qa, or ba")
		}
		if len(req.Password) < 6 {
			return fiber.NewError(400, "password must be at least 6 characters")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		var existingHash *string
		err = pool.QueryRow(context.Background(),
			`SELECT password_hash FROM employees WHERE employee_id=$1`,
			req.EmployeeID,
		).Scan(&existingHash)

		if err == nil {
			if existingHash != nil && *existingHash != "" {
				return fiber.NewError(400, "employee already registered")
			}

			_, err = pool.Exec(context.Background(), `
				UPDATE employees
				SET full_name=$2, role=$3, team=$4, password_hash=$5
				WHERE employee_id=$1
			`, req.EmployeeID, req.FullName, req.Role, req.Team, string(hash))
			if err != nil {
				return fiber.NewError(500, err.Error())
			}
		} else {
			_, err = pool.Exec(context.Background(), `
				INSERT INTO employees (employee_id, full_name, role, team, password_hash)
				VALUES ($1,$2,$3,$4,$5)
			`, req.EmployeeID, req.FullName, req.Role, req.Team, string(hash))
			if err != nil {
				return fiber.NewError(500, err.Error())
			}
		}

		token, err := createToken(req.EmployeeID, req.Role)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		user, err := getAuthUser(pool, req.EmployeeID)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		return c.JSON(AuthResponse{Token: token, User: user})
	})

	app.Post("/v1/auth/login", func(c *fiber.Ctx) error {
		var req AuthLoginReq
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(400, err.Error())
		}

		req.EmployeeID = strings.TrimSpace(req.EmployeeID)
		if req.EmployeeID == "" || req.Password == "" {
			return fiber.NewError(400, "employee_id and password are required")
		}

		var hash *string
		var role string
		err := pool.QueryRow(context.Background(),
			`SELECT password_hash, role FROM employees WHERE employee_id=$1`,
			req.EmployeeID,
		).Scan(&hash, &role)

		if err != nil || hash == nil || *hash == "" {
			return fiber.NewError(401, "invalid credentials")
		}

		if bcrypt.CompareHashAndPassword([]byte(*hash), []byte(req.Password)) != nil {
			return fiber.NewError(401, "invalid credentials")
		}

		token, err := createToken(req.EmployeeID, role)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		user, err := getAuthUser(pool, req.EmployeeID)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		return c.JSON(AuthResponse{Token: token, User: user})
	})

	app.Get("/v1/auth/me", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		out, err := getAuthUser(pool, user.EmployeeID)
		if err != nil {
			return fiber.NewError(404, "user not found")
		}

		return c.JSON(out)
	})
}
