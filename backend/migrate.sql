DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS scores;
DROP TABLE IF EXISTS managers;
DROP TABLE IF EXISTS employees;

CREATE TABLE employees (
  employee_id TEXT PRIMARY KEY,
  full_name   TEXT NOT NULL,
  role        TEXT NOT NULL CHECK (role IN ('ux','qa','ba')),
  team        TEXT NOT NULL
);

CREATE TABLE managers (
  employee_id TEXT PRIMARY KEY REFERENCES employees(employee_id)
);

CREATE TABLE submissions (
  role         TEXT NOT NULL CHECK (role IN ('ux','qa','ba')),
  employee_id  TEXT NOT NULL REFERENCES employees(employee_id),
  period_y     TEXT NOT NULL,
  submitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  payload      JSONB NOT NULL,
  PRIMARY KEY (role, employee_id, period_y)
);

CREATE TABLE scores (
  role         TEXT NOT NULL CHECK (role IN ('ux','qa','ba')),
  employee_id  TEXT NOT NULL REFERENCES employees(employee_id),
  period_y     TEXT NOT NULL,
  submitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  skill_total  DOUBLE PRECISION NOT NULL,
  credit       DOUBLE PRECISION NOT NULL,
  breakdown    JSONB NOT NULL,
  sanity_flag  TEXT NULL,
  PRIMARY KEY (role, employee_id, period_y)
);

INSERT INTO employees (employee_id, full_name, role, team) VALUES
('EMP001', 'Alice UX', 'ux', 'Design'),
('EMP002', 'Bob QA', 'qa', 'Quality'),
('EMP003', 'Carol BA', 'ba', 'Business'),
('MGR001', 'Dana Manager', 'qa', 'Management');

INSERT INTO managers (employee_id) VALUES
('MGR001');