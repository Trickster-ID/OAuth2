CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, email, password_hash, is_admin)
values ('admin', 'admin@live.com', 'JDJhJDA0JDNhRnpIL0d6d3BpNllJemdNeWlLRi5PM0NmRWZOWEM1R3dtTjl1TEpTaTlGdGJHcnB0RUpT', true);