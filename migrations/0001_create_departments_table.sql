-- +goose Up
CREATE TABLE departments (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(200) NOT NULL,
    parent_id  INT REFERENCES departments(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX unique_name_per_parent_not_null ON departments (parent_id, name) WHERE parent_id IS NOT NULL;
CREATE UNIQUE INDEX unique_name_for_root ON departments (name) WHERE parent_id IS NULL;

-- +goose Down
DROP TABLE departments;