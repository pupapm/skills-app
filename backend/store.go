package main

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScoreResponse struct {
	Role       string             `json:"role"`
	EmployeeID string             `json:"employee_id"`
	PeriodY    string             `json:"period_y"`
	SkillTotal float64            `json:"skill_total"`
	Credit     float64            `json:"credit"`
	Breakdown  map[string]float64 `json:"breakdown"`
	SanityFlag *string            `json:"sanity_flag"`
}

func upsertSubmissionAndScore(
	c *fiber.Ctx,
	pool *pgxpool.Pool,
	role string,
	employeeID string,
	periodY string,
	payload any,
	skillTotal float64,
	credit float64,
	breakdown map[string]float64,
	sanity *string,
) (*ScoreResponse, error) {
	ctx := context.Background()

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	breakdownJSON, err := json.Marshal(breakdown)
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO submissions (role, employee_id, period_y, payload)
		VALUES ($1, $2, $3, $4::jsonb)
		ON CONFLICT (role, employee_id, period_y)
		DO UPDATE SET
			payload = EXCLUDED.payload,
			submitted_at = now()
	`, role, employeeID, periodY, string(payloadJSON))
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO scores (role, employee_id, period_y, skill_total, credit, breakdown, sanity_flag)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7)
		ON CONFLICT (role, employee_id, period_y)
		DO UPDATE SET
			skill_total = EXCLUDED.skill_total,
			credit = EXCLUDED.credit,
			breakdown = EXCLUDED.breakdown,
			sanity_flag = EXCLUDED.sanity_flag,
			submitted_at = now()
	`, role, employeeID, periodY, skillTotal, credit, string(breakdownJSON), sanity)
	if err != nil {
		return nil, err
	}

	return &ScoreResponse{
		Role:       role,
		EmployeeID: employeeID,
		PeriodY:    periodY,
		SkillTotal: skillTotal,
		Credit:     credit,
		Breakdown:  breakdown,
		SanityFlag: sanity,
	}, nil
}
