package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeProject struct {
	ID             int64  `json:"id"`
	EmployeeID     string `json:"employee_id"`
	ProjectName    string `json:"project_name"`
	ProjectYear    string `json:"project_year"`
	Responsibility string `json:"responsibility"`
	WhatYouDid     string `json:"what_you_did"`
	ToolsUsed      string `json:"tools_used"`
	Outcome        string `json:"outcome"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type EmployeeProjectPayload struct {
	ProjectName    string `json:"project_name"`
	ProjectYear    string `json:"project_year"`
	Responsibility string `json:"responsibility"`
	WhatYouDid     string `json:"what_you_did"`
	ToolsUsed      string `json:"tools_used"`
	Outcome        string `json:"outcome"`
}

func validateProjectPayload(p EmployeeProjectPayload) error {
	if strings.TrimSpace(p.ProjectName) == "" {
		return fiber.NewError(400, "project_name is required")
	}
	if strings.TrimSpace(p.ProjectYear) == "" {
		return fiber.NewError(400, "project_year is required")
	}
	return nil
}

func scanProjectRows(rows pgxRowsLike) ([]EmployeeProject, error) {
	var out []EmployeeProject
	for rows.Next() {
		var p EmployeeProject
		if err := rows.Scan(
			&p.ID, &p.EmployeeID, &p.ProjectName, &p.ProjectYear,
			&p.Responsibility, &p.WhatYouDid, &p.ToolsUsed, &p.Outcome,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

type pgxRowsLike interface {
	Next() bool
	Scan(dest ...any) error
}

func registerProjectRoutes(app *fiber.App, pool *pgxpool.Pool) {
	app.Get("/v1/projects", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		rows, err := pool.Query(context.Background(), `
			SELECT id, employee_id, project_name, project_year, responsibility, what_you_did, tools_used, outcome,
			       created_at::text, updated_at::text
			FROM employee_projects
			WHERE employee_id = $1
			ORDER BY project_year DESC, updated_at DESC, id DESC
		`, user.EmployeeID)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}
		defer rows.Close()

		out, err := scanProjectRows(rows)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		return c.JSON(out)
	})

	app.Post("/v1/projects", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		var payload EmployeeProjectPayload
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(400, err.Error())
		}
		if err := validateProjectPayload(payload); err != nil {
			return err
		}

		var p EmployeeProject
		err := pool.QueryRow(context.Background(), `
			INSERT INTO employee_projects (
				employee_id, project_name, project_year, responsibility, what_you_did, tools_used, outcome
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			RETURNING id, employee_id, project_name, project_year, responsibility, what_you_did, tools_used, outcome,
			          created_at::text, updated_at::text
		`,
			user.EmployeeID,
			strings.TrimSpace(payload.ProjectName),
			strings.TrimSpace(payload.ProjectYear),
			strings.TrimSpace(payload.Responsibility),
			strings.TrimSpace(payload.WhatYouDid),
			strings.TrimSpace(payload.ToolsUsed),
			strings.TrimSpace(payload.Outcome),
		).Scan(
			&p.ID, &p.EmployeeID, &p.ProjectName, &p.ProjectYear,
			&p.Responsibility, &p.WhatYouDid, &p.ToolsUsed, &p.Outcome,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		return c.JSON(p)
	})

	app.Put("/v1/projects/:id", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return fiber.NewError(400, "invalid project id")
		}

		var payload EmployeeProjectPayload
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(400, err.Error())
		}
		if err := validateProjectPayload(payload); err != nil {
			return err
		}

		var p EmployeeProject
		err = pool.QueryRow(context.Background(), `
			UPDATE employee_projects
			SET project_name = $1,
			    project_year = $2,
			    responsibility = $3,
			    what_you_did = $4,
			    tools_used = $5,
			    outcome = $6,
			    updated_at = now()
			WHERE id = $7 AND employee_id = $8
			RETURNING id, employee_id, project_name, project_year, responsibility, what_you_did, tools_used, outcome,
			          created_at::text, updated_at::text
		`,
			strings.TrimSpace(payload.ProjectName),
			strings.TrimSpace(payload.ProjectYear),
			strings.TrimSpace(payload.Responsibility),
			strings.TrimSpace(payload.WhatYouDid),
			strings.TrimSpace(payload.ToolsUsed),
			strings.TrimSpace(payload.Outcome),
			id,
			user.EmployeeID,
		).Scan(
			&p.ID, &p.EmployeeID, &p.ProjectName, &p.ProjectYear,
			&p.Responsibility, &p.WhatYouDid, &p.ToolsUsed, &p.Outcome,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return fiber.NewError(404, "project not found")
		}

		return c.JSON(p)
	})

	app.Delete("/v1/projects/:id", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return fiber.NewError(400, "invalid project id")
		}

		tag, err := pool.Exec(context.Background(), `
			DELETE FROM employee_projects
			WHERE id = $1 AND employee_id = $2
		`, id, user.EmployeeID)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}
		if tag.RowsAffected() == 0 {
			return fiber.NewError(404, "project not found")
		}

		return c.JSON(fiber.Map{"ok": true})
	})

	app.Get("/v1/manager/employee/projects", AuthRequired(), ManagerOnly(pool), func(c *fiber.Ctx) error {
		employeeID := strings.TrimSpace(c.Query("employee_id"))
		if employeeID == "" {
			return fiber.NewError(400, "employee_id is required")
		}

		rows, err := pool.Query(context.Background(), `
			SELECT id, employee_id, project_name, project_year, responsibility, what_you_did, tools_used, outcome,
			       created_at::text, updated_at::text
			FROM employee_projects
			WHERE employee_id = $1
			ORDER BY project_year DESC, updated_at DESC, id DESC
		`, employeeID)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}
		defer rows.Close()

		out, err := scanProjectRows(rows)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		return c.JSON(out)
	})
}
