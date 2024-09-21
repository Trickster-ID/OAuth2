CREATE TABLE roles (
    id INT GENERATED ALWAYS AS IDENTITY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);

INSERT INTO roles (code, name, description, is_active)
VALUES
('admin', 'Administrator', 'as administrator, able to access anything on app with this auth app', true),
('user', 'Users', 'common user role with limited access', true);

CREATE TABLE users (
    id INT GENERATED ALWAYS AS IDENTITY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    CONSTRAINT fk_role
        FOREIGN KEY(role_id)
            REFERENCES roles(id)
);

INSERT INTO users (username, email, password_hash, role_id)
values ('admin', 'admin@live.com', '$2a$04$ZP1.DVAdR677eHTBUpDzE.0hHnp31JcyRK/eMF9Z7Y.iOWJkE/JNi', 1);

