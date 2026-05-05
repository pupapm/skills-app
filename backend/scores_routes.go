package main

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func registerScoreRoutes(app *fiber.App, pool *pgxpool.Pool) {
	app.Get("/v1/scores/latest", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		periodY := c.Query("period_y")
		if periodY == "" {
			return fiber.NewError(400, "period_y is required")
		}

		var role, employeeID, fetchedPeriodY string
		var skillTotal, credit float64
		var breakdownJSON []byte
		var sanity *string

		err := pool.QueryRow(context.Background(), `
			SELECT role, employee_id, period_y, skill_total, credit, breakdown, sanity_flag
			FROM scores
			WHERE role = $1 AND employee_id = $2 AND period_y = $3
		`, user.Role, user.EmployeeID, periodY).Scan(
			&role, &employeeID, &fetchedPeriodY, &skillTotal, &credit, &breakdownJSON, &sanity,
		)
		if err != nil {
			return fiber.NewError(404, "score not found")
		}

		breakdown := map[string]float64{}
		_ = json.Unmarshal(breakdownJSON, &breakdown)

		return c.JSON(ScoreResponse{
			Role:       role,
			EmployeeID: employeeID,
			PeriodY:    fetchedPeriodY,
			SkillTotal: skillTotal,
			Credit:     credit,
			Breakdown:  breakdown,
			SanityFlag: sanity,
		})
	})

	app.Get("/v1/scores/history", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)
		limit := c.QueryInt("limit", 12)

		rows, err := pool.Query(context.Background(), `
			SELECT role, employee_id, period_y, skill_total, credit, breakdown, sanity_flag
			FROM scores
			WHERE role = $1 AND employee_id = $2
			ORDER BY period_y DESC
			LIMIT $3
		`, user.Role, user.EmployeeID, limit)
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
