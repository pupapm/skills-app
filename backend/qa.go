package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QASubmit struct {
	PeriodY string `json:"period_y"`

	TenureBucket string `json:"qa_tenure_bucket"`

	MainProjects    int `json:"qa_main_projects"`
	SupportProjects int `json:"qa_support_projects"`

	BugsReportedTotal     int `json:"bugs_reported_total"`
	CriticalIncidentCount int `json:"critical_incident_count"`

	FixBucket    string `json:"fix_no_question_ratio_bucket"`
	ReopenBucket string `json:"reopen_due_info_ratio_bucket"`

	BugWithEvidenceCount    int `json:"bug_with_evidence_count"`
	AutomationExecutedCount int `json:"automation_executed_count"`
}

func bucket01(s string) float64 {
	switch s {
	case "0–20%", "0-20%":
		return 0.2
	case "21–40%", "21-40%":
		return 0.4
	case "41–60%", "41-60%":
		return 0.6
	case "61–80%", "61-80%":
		return 0.8
	case "81–100%", "81-100%":
		return 1.0
	default:
		return 0
	}
}

func qaTimeFactor(bucket string) float64 {
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

func scoreQA(p QASubmit) (skill float64, credit float64, breakdown map[string]float64, sanity *string) {
	breakdown = map[string]float64{}

	if p.MainProjects < 0 || p.SupportProjects < 0 ||
		p.BugsReportedTotal < 0 || p.CriticalIncidentCount < 0 ||
		p.BugWithEvidenceCount < 0 || p.AutomationExecutedCount < 0 {
		s := "NEGATIVE_INPUT"
		sanity = &s
	}

	if p.BugsReportedTotal > 0 && p.BugWithEvidenceCount > p.BugsReportedTotal {
		s := "EVIDENCE_GT_BUGS"
		sanity = &s
	}

	exp := qaTimeFactor(p.TenureBucket) * 15
	mainDen := float64(maxInt(p.MainProjects, 1))

	incidentRate := float64(p.CriticalIncidentCount) / mainDen
	risk := (1 - cap01(incidentRate/0.20)) * 35

	fixR := bucket01(p.FixBucket)
	reopenR := bucket01(p.ReopenBucket)
	comm := cap01(fixR-0.5*reopenR) * 20

	bugDen := float64(maxInt(p.BugsReportedTotal, 1))
	evidenceRatio := float64(p.BugWithEvidenceCount) / bugDen
	evidence := cap01(evidenceRatio/0.70) * 20
	auto := cap01(float64(p.AutomationExecutedCount)/10.0) * 10
	tool := evidence + auto

	base := exp + risk + comm + tool
	credit = cap01(float64(p.SupportProjects)/20.0) * 10
	final := base + credit

	skill = capRange(final, 0, 110)

	breakdown["exp"] = exp
	breakdown["risk"] = risk
	breakdown["comm"] = comm
	breakdown["tooling"] = tool
	breakdown["base"] = base
	breakdown["credit"] = credit
	breakdown["final"] = skill

	return
}

func registerQA(app *fiber.App, pool *pgxpool.Pool) {
	app.Post("/v1/qa/submit", AuthRequired(), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*User)

		if user.Role != "qa" {
			return fiber.NewError(403, "role mismatch")
		}

		var p QASubmit
		if err := c.BodyParser(&p); err != nil {
			return fiber.NewError(400, err.Error())
		}
		if p.PeriodY == "" {
			return fiber.NewError(400, "period_y is required")
		}

		skill, credit, breakdown, sanity := scoreQA(p)

		out, err := upsertSubmissionAndScore(
			c,
			pool,
			"qa",
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
