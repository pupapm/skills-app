package main

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ManagerMemberRow struct {
	EmployeeID string  `json:"employee_id"`
	FullName   string  `json:"full_name"`
	Team       string  `json:"team"`
	SkillTotal float64 `json:"skill_total"`
	Credit     float64 `json:"credit"`
	SanityFlag *string `json:"sanity_flag"`
}

func registerManagerRoutes(app *fiber.App, pool *pgxpool.Pool) {
	app.Get("/v1/manager/members", AuthRequired(), ManagerOnly(pool), func(c *fiber.Ctx) error {
		periodY := c.Query("period_y")
		role := c.Query("role")
		if periodY == "" {
			return fiber.NewError(400, "period_y is required")
		}
		if role == "" {
			return fiber.NewError(400, "role is required")
		}

		rows, err := pool.Query(context.Background(), `
			SELECT e.employee_id, e.full_name, e.team,
			       COALESCE(s.skill_total, 0),
			       COALESCE(s.credit, 0),
			       s.sanity_flag
			FROM employees e
			LEFT JOIN scores s
			  ON s.employee_id = e.employee_id
			 AND s.role = e.role
			 AND s.period_y = $1
			WHERE e.role = $2
			ORDER BY e.employee_id
		`, periodY, role)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}
		defer rows.Close()

		var out []ManagerMemberRow
		for rows.Next() {
			var r ManagerMemberRow
			if err := rows.Scan(&r.EmployeeID, &r.FullName, &r.Team, &r.SkillTotal, &r.Credit, &r.SanityFlag); err != nil {
				return fiber.NewError(500, err.Error())
			}
			out = append(out, r)
		}
		return c.JSON(out)
	})

	app.Get("/v1/manager/employee", AuthRequired(), ManagerOnly(pool), func(c *fiber.Ctx) error {
		employeeID := c.Query("employee_id")
		periodY := c.Query("period_y")
		if employeeID == "" {
			return fiber.NewError(400, "employee_id is required")
		}
		if periodY == "" {
			return fiber.NewError(400, "period_y is required")
		}

		var role, fullName, team, fetchedEmployeeID, fetchedPeriodY string
		var skillTotal, credit float64
		var breakdownJSON []byte
		var sanity *string

		err := pool.QueryRow(context.Background(), `
			SELECT e.role, e.full_name, e.team, s.employee_id, s.period_y, s.skill_total, s.credit, s.breakdown, s.sanity_flag
			FROM scores s
			JOIN employees e
			  ON e.employee_id = s.employee_id AND e.role = s.role
			WHERE s.employee_id = $1 AND s.period_y = $2
		`, employeeID, periodY).Scan(
			&role, &fullName, &team, &fetchedEmployeeID, &fetchedPeriodY, &skillTotal, &credit, &breakdownJSON, &sanity,
		)
		if err != nil {
			return fiber.NewError(404, "employee score not found")
		}

		breakdown := map[string]float64{}
		_ = json.Unmarshal(breakdownJSON, &breakdown)

		return c.JSON(fiber.Map{
			"role":        role,
			"employee_id": fetchedEmployeeID,
			"full_name":   fullName,
			"team":        team,
			"period_y":    fetchedPeriodY,
			"skill_total": skillTotal,
			"credit":      credit,
			"breakdown":   breakdown,
			"sanity_flag": sanity,
		})
	})

	app.Get("/v1/manager/employee/history", AuthRequired(), ManagerOnly(pool), func(c *fiber.Ctx) error {
		employeeID := c.Query("employee_id")
		limit := c.QueryInt("limit", 12)
		if employeeID == "" {
			return fiber.NewError(400, "employee_id is required")
		}

		rows, err := pool.Query(context.Background(), `
			SELECT role, employee_id, period_y, skill_total, credit, breakdown, sanity_flag
			FROM scores
			WHERE employee_id = $1
			ORDER BY period_y DESC
			LIMIT $2
		`, employeeID, limit)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}
		defer rows.Close()

		var out []ScoreResponse
		for rows.Next() {
			var r ScoreResponse
			var breakdownJSON []byte
			if err := rows.Scan(
				&r.Role, &r.EmployeeID, &r.PeriodY, &r.SkillTotal, &r.Credit, &breakdownJSON, &r.SanityFlag,
			); err != nil {
				return fiber.NewError(500, err.Error())
			}
			_ = json.Unmarshal(breakdownJSON, &r.Breakdown)
			out = append(out, r)
		}
		return c.JSON(out)
	})
}
