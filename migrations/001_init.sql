-- +goose Up

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE departments (
  id          UUID PRIMARY KEY,
  name        VARCHAR(200) NOT NULL,
  parent_id   UUID NULL REFERENCES departments(id) ON DELETE CASCADE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_departments_parent_name ON departments(parent_id, name);
CREATE UNIQUE INDEX ux_departments_root_name ON departments(name) WHERE parent_id IS NULL;
CREATE INDEX idx_departments_parent_id ON departments(parent_id);

CREATE TABLE employees (
  id            UUID PRIMARY KEY,
  department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
  full_name     VARCHAR(200) NOT NULL,
  position      VARCHAR(200) NOT NULL,
  hired_at      DATE NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_employees_department_id ON employees(department_id);

-- +goose Down
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;
DROP EXTENSION IF EXISTS pgcrypto;