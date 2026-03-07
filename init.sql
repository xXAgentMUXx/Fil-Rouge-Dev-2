CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(20) NOT NULL DEFAULT 'client'
);

CREATE TABLE agencies (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    city VARCHAR(100)
);

CREATE TABLE properties (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    city VARCHAR(100),
    price NUMERIC,
    surface INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    agency_id UUID,
    is_sold BOOLEAN DEFAULT false,
    agent_id INTEGER,
    image VARCHAR(255),

    CONSTRAINT properties_agency_id_fkey
        FOREIGN KEY (agency_id) REFERENCES agencies(id),

    CONSTRAINT properties_agent_id_fkey
        FOREIGN KEY (agent_id) REFERENCES users(id)
);

CREATE TABLE favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    property_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_favorite UNIQUE (user_id, property_id),

    CONSTRAINT favorites_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT favorites_property_id_fkey
        FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    property_id INTEGER,
    buyer_id INTEGER,
    amount NUMERIC,
    stripe_session_id TEXT,
    status TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT payments_property_id_fkey
        FOREIGN KEY (property_id) REFERENCES properties(id),

    CONSTRAINT payments_buyer_id_fkey
        FOREIGN KEY (buyer_id) REFERENCES users(id)
);

CREATE TABLE sales (
    id SERIAL PRIMARY KEY,
    property_id INTEGER NOT NULL,
    buyer_id INTEGER NOT NULL,
    sale_price NUMERIC NOT NULL,
    sold_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT sales_property_id_fkey
        FOREIGN KEY (property_id) REFERENCES properties(id),

    CONSTRAINT sales_buyer_id_fkey
        FOREIGN KEY (buyer_id) REFERENCES users(id)
);