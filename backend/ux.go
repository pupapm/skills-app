package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UXSubmit struct {
	PeriodY string `json:"period_y"`

	TenureBucket string `json:"ux_tenure_bucket"`

	MainProjects    int `json:"ux_main_projects"`
	SupportProjects int `json:"ux_support_projects"`

	ReworkDueScopeMain          int `json:"rework_due_scope_main"`
	SolutionsProposedMain       int `json:"solutions_proposed_main"`
	LateEntryMain               int `json:"ux_late_entry_main"`
	UsabilityIssuesAfterRelease int `json:"usability_issues_after_release"`
}

func uxTimeFactor(bucket string) float64 {
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

func scoreUX(p UXSubmit) (skill float64, credit float64, breakdown map[string]float64, sanity *string) {
	breakdown = map[string]float64{}

	if p.MainProjects < 0 || p.SupportProjects < 0 ||
		p.ReworkDueScopeMain < 0 || p.SolutionsProposedMain < 0 ||
		p.LateEntryMain < 0 || p.UsabilityIssuesAfterRelease < 0 {
		s := "NEGATIVE_INPUT"
		sanity = &s
	}

	if p.MainProjects > 0 {
		if p.ReworkDueScopeMain > p.MainProjects {
			s := "OVERCOUNT_REWORK"
			sanity = &s
		}
		if p.SolutionsProposedMain > p.MainProjects {
			s := "OVERCOUNT_OPTIONS"
			sanity = &s
		}
		if p.LateEntryMain > p.MainProjects {
			s := "OVERCOUNT_LATE"
			sanity = &s
		}
		if p.UsabilityIssuesAfterRelease > p.MainProjects {
			s := "OVERCOUNT_USABILITY_ISSUES"
			sanity = &s
		}
	}

	exp := uxTimeFactor(p.TenureBucket) * 15
	den := float64(maxInt(p.MainProjects, 1))

	reworkRate := float64(p.ReworkDueScopeMain) / den
	scope := (1 - cap01(reworkRate/0.30)) * 25

	optRate := float64(p.SolutionsProposedMain) / den
	solution := cap01(optRate/0.50) * 25

	lateRate := float64(p.LateEntryMain) / den
	delivery := (1 - cap01(lateRate/0.20)) * 20

	issuesRate := float64(p.UsabilityIssuesAfterRelease) / den
	impact := (1 - cap01(issuesRate/0.30)) * 15

	base := exp + scope + solution + delivery + impact
	credit = cap01(float64(p.SupportProjects)/20.0) * 10
	final := base + credit

	skill = capRange(final, 0, 110)

	breakdown["exp"] = exp
	breakdown["scope"] = scope
	breakdown["solution"] = solution
	breakdown["delivery"] = delivery
	breakdown["impact"] = impact
	breakdown["base"] = base
	breakdown["credit"] = credit
	breakdown["final"] = skill

	return
}

func registerUX(app *fiber.App, pool *pgxpool.Pool) {
	app.Post("/v1/ux/submit", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		if user.Role != "ux" {
			return fiber.NewError(403, "role mismatch")
		}

		var p UXSubmit
		if err := c.BodyParser(&p); err != nil {
			return fiber.NewError(400, err.Error())
		}
		if p.PeriodY == "" {
			return fiber.NewError(400, "period_y is required")
		}

		skill, credit, breakdown, sanity := scoreUX(p)

		out, err := upsertSubmissionAndScore(
			c,
			pool,
			"ux",
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
