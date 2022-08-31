DROP TYPE IF EXISTS dim_users;

 CREATE TABLE dim_users (
    id SERIAL PRIMARY KEY,
    created_dt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(200) DEFAULT 'Anonymus'
);