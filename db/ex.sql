CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE agencies (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    city VARCHAR(100)
);

CREATE TABLE properties (
    id UUID PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    city VARCHAR(100),
    price DECIMAL,
    surface INT,
    agency_id UUID REFERENCES agencies(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sales (
    id UUID PRIMARY KEY,
    property_id UUID REFERENCES properties(id),
    buyer_id UUID REFERENCES users(id),
    sale_price DECIMAL,
    sold_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);