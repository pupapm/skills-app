package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BASubmit struct {
	PeriodY string `json:"period_y"`

	TenureBucket string `json:"ba_tenure_bucket"`

	MainProjects    int `json:"ba_main_projects"`
	SupportProjects int `json:"ba_support_projects"`

	ReqReworkMain      int `json:"req_rework_due_ambiguity_main"`
	ActionableNoRepeat int `json:"actionable_req_no_repeat_main"`
	SignoffBeforeBuild int `json:"decision_signoff_before_build_main"`
	RiskPreventedCount int `json:"risk_prevented_count_main"`
}

func baTimeFactor(bucket string) float64 {
	switch bucket {
	case "< 6 เดือน":
		return 0.2
	case "6 - 12 เดือน":
		return 0.4
	case "1 - 2 ปี":
		return 0.6
	case "2 - 4 ปี":
		return 0.8
	case "> 4 ปี":
		return 1.0
	default:
		return 0
	}
}

func scoreBA(p BASubmit) (skill float64, credit float64, breakdown map[string]float64, sanity *string) {
	breakdown = map[string]float64{}

	if p.MainProjects < 0 || p.SupportProjects < 0 ||
		p.ReqReworkMain < 0 || p.ActionableNoRepeat < 0 ||
		p.SignoffBeforeBuild < 0 || p.RiskPreventedCount < 0 {
		s := "NEGATIVE_INPUT"
		sanity = &s
	}

	if p.MainProjects > 0 {
		if p.ReqReworkMain > p.MainProjects {
			s := "OVERCOUNT_REWORK"
			sanity = &s
		}
		if p.ActionableNoRepeat > p.MainProjects {
			s := "OVERCOUNT_ACTIONABLE"
			sanity = &s
		}
		if p.SignoffBeforeBuild > p.MainProjects {
			s := "OVERCOUNT_SIGNOFF"
			sanity = &s
		}
	}

	exp := baTimeFactor(p.TenureBucket) * 15
	den := float64(maxInt(p.MainProjects, 1))

	reworkRate := float64(p.ReqReworkMain) / den
	scope := (1 - cap01(reworkRate/0.30)) * 30

	actionableRate := float64(p.ActionableNoRepeat) / den
	signoffRate := float64(p.SignoffBeforeBuild) / den
	decision := (0.6*cap01(actionableRate/0.70) + 0.4*cap01(signoffRate/0.80)) * 35

	riskRate := float64(p.RiskPreventedCount) / den
	risk := cap01(riskRate/0.30) * 20

	base := exp + scope + decision + risk
	credit = cap01(float64(p.SupportProjects)/20.0) * 10
	final := base + credit

	skill = capRange(final, 0, 110)

	breakdown["exp"] = exp
	breakdown["scope"] = scope
	breakdown["decision"] = decision
	breakdown["risk"] = risk
	breakdown["base"] = base
	breakdown["credit"] = credit
	breakdown["final"] = skill

	return
}

func registerBA(app *fiber.App, pool *pgxpool.Pool) {
	app.Post("/v1/ba/submit", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		if user.Role != "ba" {
			return fiber.NewError(403, "role mismatch")
		}

		var p BASubmit
		if err := c.BodyParser(&p); err != nil {
			return fiber.NewError(400, err.Error())
		}
		if p.PeriodY == "" {
			return fiber.NewError(400, "period_y is required")
		}

		skill, credit, breakdown, sanity := scoreBA(p)

		out, err := upsertSubmissionAndScore(
			c,
			pool,
			"ba",
			user.EmployeeID,
			p.PeriodY,
			p,
			skill,
			credit,
			breakdown,
			sanity,
		)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		return c.JSON(out)
	})
}
