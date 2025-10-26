CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50) UNIQUE,
    address TEXT,
    role VARCHAR(50) DEFAULT 'CUSTOMER',
    status VARCHAR(50) DEFAULT 'ACTIVE',
    subscription_status VARCHAR(50) DEFAULT 'SUBSCRIBED',
    password TEXT NOT NULL,
    token TEXT,
    refresh_token TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);
