export function parseBreakdown(input) {
  if (!input) return {};
  if (typeof input === "string") {
    try {
      return JSON.parse(input);
    } catch {
      return {};
    }
  }
  if (typeof input === "object") return input;
  return {};
}

export function formatRoleBreakdown(role, breakdown) {
  const b = parseBreakdown(breakdown);

  const config = {
    ux: [
      { key: "exp", label: "Experience", max: 15 },
      { key: "scope", label: "Scope Stability", max: 25 },
      { key: "solution", label: "Solution Quality", max: 25 },
      { key: "delivery", label: "Delivery Discipline", max: 20 },
      { key: "impact", label: "Usability Impact", max: 15 },
      { key: "credit", label: "Contribution Credit", max: 10 },
    ],

    qa: [
      { key: "exp", label: "Experience", max: 15 },
      { key: "risk", label: "Risk Control", max: 35 },
      { key: "comm", label: "Communication", max: 20 },
      { key: "tooling", label: "Evidence & Automation", max: 30 },
      { key: "credit", label: "Contribution Credit", max: 10 },
    ],

    ba: [
      { key: "exp", label: "Experience", max: 15 },
      { key: "scope", label: "Requirement Stability", max: 30 },
      { key: "decision", label: "Decision Quality", max: 35 },
      { key: "risk", label: "Risk Prevention", max: 20 },
      { key: "credit", label: "Contribution Credit", max: 10 },
    ],

    dev: [
        { key: "exp", label: "Experience", max: 15 },
        { key: "quality", label: "Code Quality", max: 30 },
        { key: "delivery", label: "Delivery Discipline", max: 25 },
        { key: "engineering", label: "Engineering Practice", max: 20 },
        { key: "credit", label: "Contribution Credit", max: 10 },
    ],
  };

  const selected = config[role] || Object.keys(b).map((key) => ({
    key,
    label: key,
    max: 100,
  }));

  return selected
    .filter((item) => item.key in b)
    .map((item) => ({
      key: item.key,
      label: item.label,
      value: Number(b[item.key] || 0),
      max: item.max,
      percent: item.max > 0 ? (Number(b[item.key] || 0) / item.max) * 100 : 0,
    }));
}

export function avg(arr, key) {
  if (!arr.length) return 0;
  return arr.reduce((sum, x) => sum + Number(x[key] || 0), 0) / arr.length;
}

export function maxValue(arr, key) {
  if (!arr.length) return 0;
  return Math.max(...arr.map((x) => Number(x[key] || 0)));
}

export function minValue(arr, key) {
  if (!arr.length) return 0;
  return Math.min(...arr.map((x) => Number(x[key] || 0)));
}