CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS places (
    id            uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    risk_type_id  uuid NOT NULL,
    place_type_id uuid NOT NULL,
    latitude      DECIMAL(10,6) NOT NULL,
    longitude     DECIMAL(10,6) NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NULL,
    deleted_at    TIMESTAMP NULL
);