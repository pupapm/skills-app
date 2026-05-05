package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DevSubmit struct {
	PeriodY string `json:"period_y"`

	TenureBucket string `json:"dev_tenure_bucket"`

	MainProjects    int `json:"dev_main_projects"`
	SupportProjects int `json:"dev_support_projects"`

	BugFromOwnCode int `json:"bug_from_own_code"`

	CleanImplementation int `json:"clean_implementation"`

	OnTimeDelivery int `json:"on_time_delivery"`
	BlockerCaused  int `json:"blocker_caused"`

	CodeReviewContribution int `json:"code_review_contribution"`
	OptimizationWork       int `json:"optimization_work"`
}

func devTimeFactor(bucket string) float64 {
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

func scoreDev(p DevSubmit) (skill float64, credit float64, breakdown map[string]float64, sanity *string) {
	breakdown = map[string]float64{}

	// sanity
	if p.MainProjects < 0 || p.SupportProjects < 0 {
		s := "NEGATIVE_INPUT"
		sanity = &s
	}

	if p.MainProjects > 0 {
		if p.BugFromOwnCode > p.MainProjects {
			s := "OVERCOUNT_BUG"
			sanity = &s
		}
		if p.CleanImplementation > p.MainProjects {
			s := "OVERCOUNT_CLEAN"
			sanity = &s
		}
		if p.OnTimeDelivery > p.MainProjects {
			s := "OVERCOUNT_DELIVERY"
			sanity = &s
		}
		if p.BlockerCaused > p.MainProjects {
			s := "OVERCOUNT_BLOCKER"
			sanity = &s
		}
	}

	// Experience (15)
	exp := devTimeFactor(p.TenureBucket) * 15
	breakdown["exp"] = exp

	den := float64(maxInt(p.MainProjects, 1))

	// 🔹 Code Quality (30)
	bugRate := float64(p.BugFromOwnCode) / den
	bugPenalty := cap01(bugRate / 0.30)

	cleanRate := float64(p.CleanImplementation) / den
	cleanScore := cap01(cleanRate / 0.70)

	quality := (0.5*(1-bugPenalty) + 0.5*cleanScore) * 30
	breakdown["quality"] = quality

	// 🔹 Delivery (25)
	onTimeRate := float64(p.OnTimeDelivery) / den
	blockerRate := float64(p.BlockerCaused) / den

	delivery := (0.6*cap01(onTimeRate/0.80) + 0.4*(1-cap01(blockerRate/0.20))) * 25
	breakdown["delivery"] = delivery

	// 🔹 Engineering (20)
	reviewScore := cap01(float64(p.CodeReviewContribution) / 20.0)
	optScore := cap01(float64(p.OptimizationWork) / 10.0)

	engineer := (0.5*reviewScore + 0.5*optScore) * 20
	breakdown["engineering"] = engineer

	// 🔹 Base
	base := exp + quality + delivery + engineer
	breakdown["base"] = base

	// 🔹 Credit (10)
	credit = cap01(float64(p.SupportProjects)/20.0) * 10
	breakdown["credit"] = credit

	// 🔹 Final
	final := base + credit
	skill = capRange(final, 0, 110)
	breakdown["final"] = skill

	return
}

func registerDev(app *fiber.App, pool any) {
	app.Post("/v1/dev/submit", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		if user.Role != "dev" {
			return fiber.NewError(403, "role mismatch")
		}

		var p DevSubmit
		if err := c.BodyParser(&p); err != nil {
			return fiber.NewError(400, err.Error())
		}

		skill, credit, breakdown, sanity := scoreDev(p)

		pg := pool.(*pgxpool.Pool)

		out, err := upsertSubmissionAndScore(
			c,
			pg,
			"dev",
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
