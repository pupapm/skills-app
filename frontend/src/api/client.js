export function setToken(token) {
  localStorage.setItem("token", token);
}

export function getToken() {
  return localStorage.getItem("token") || "";
}

export function setSession(user) {
  localStorage.setItem("session", JSON.stringify(user || {}));
}

export function getSession() {
  try {
    return JSON.parse(localStorage.getItem("session") || "{}");
  } catch {
    return {};
  }
}

export function clearSession() {
  localStorage.removeItem("token");
  localStorage.removeItem("session");
}

async function request(path, options = {}) {
  const token = getToken();

  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const API_BASE = import.meta.env.VITE_API_BASE_URL || "";

  const res = await fetch(API_BASE + path, {
    ...options,
    headers,
  });

  const data = await res.json().catch(() => ({}));

  if (!res.ok) {
    throw new Error(data?.message || data?.error || "Request failed");
  }

  return data;
}

export const api = {
  register: (payload) =>
    request("/v1/auth/register", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  login: (payload) =>
    request("/v1/auth/login", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  me: () => request("/v1/auth/me"),

  latestScore: (periodY) =>
    request(`/v1/scores/latest?period_y=${encodeURIComponent(periodY)}`),

  history: (limit = 12) =>
    request(`/v1/scores/history?limit=${limit}`),

  submitUX: (payload) =>
    request("/v1/ux/submit", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  submitQA: (payload) =>
    request("/v1/qa/submit", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  submitBA: (payload) =>
    request("/v1/ba/submit", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  managerMembers: (periodY, role) =>
    request(
      `/v1/manager/members?period_y=${encodeURIComponent(periodY)}&role=${encodeURIComponent(role)}`
    ),

  managerEmployee: (employeeId, periodY) =>
    request(
      `/v1/manager/employee?employee_id=${encodeURIComponent(employeeId)}&period_y=${encodeURIComponent(periodY)}`
    ),

  managerEmployeeHistory: (employeeId, limit = 12) =>
    request(
      `/v1/manager/employee/history?employee_id=${encodeURIComponent(employeeId)}&limit=${limit}`
    ),

  myProjects: () => request("/v1/projects"),

  createProject: (payload) =>
    request("/v1/projects", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  updateProject: (id, payload) =>
    request(`/v1/projects/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
    }),

  deleteProject: (id) =>
    request(`/v1/projects/${id}`, {
      method: "DELETE",
    }),

  managerEmployeeProjects: (employeeId) =>
    request(
      `/v1/manager/employee/projects?employee_id=${encodeURIComponent(employeeId)}`
    ),
};