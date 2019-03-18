
CREATE TABLE `user` (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    uuid VARCHAR(36) NOT NULL,
    first_name VARCHAR(24) NULL,
    last_name VARCHAR(24) NULL,
    email VARCHAR(128) NULL,
    phone VARCHAR(24) NOT NULL,
    status_id TINYINT(3) UNSIGNED NOT NULL DEFAULT 1,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `user` ADD CONSTRAINT uk_user_email UNIQUE(email);

CREATE TABLE `ref_user_status` (
    id TINYINT(3) UNSIGNED NOT NULL PRIMARY KEY,
    name VARCHAR(16) NOT NULL,
    slug VARCHAR(16) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `ref_user_status` (id, name, slug)
VALUES
(1, "Unverified", "unverified"),
(2, "Verified", "verified"),
(3, "Deleted", "deleted");

ALTER TABLE `user` ADD FOREIGN KEY (status_id) REFERENCES ref_user_status(id);

-- Insert dummy user
INSERT INTO user (uuid, first_name, last_name, email, phone, status_id, hashed_password, created)
VALUES (
    UUID(),
    'Alice',
    'Jones',
    'alice@example.com',
    '867-5309',
    2,
    '$2a$12$3/ZmSDnMwcfRcxgahkkNrOzyOv28HXbEu1vSVdkIbFId.AKJryrDC', -- 'password'
    '2019-02-24 17:25:25'
);

-- CREATE store table
CREATE TABLE `store` (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    phone VARCHAR(32) NULL,
    email VARCHAR(32) NULL,
    url VARCHAR(64) NULL,
    address VARCHAR(128) NULL,
    city VARCHAR(64) NULL,
    state VARCHAR(32) NULL,
    zip VARCHAR(16) NULL,
    lat DECIMAL(9, 6),
    lng DECIMAL(9, 6),
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create ingredient table
CREATE TABLE `ingredient` (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create flavor table
CREATE TABLE `flavor` (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    description TEXT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create flavor_ingredient table
CREATE TABLE `flavor_ingredient` (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    flavor_id INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `flavor_ingredient` ADD FOREIGN KEY (flavor_id) REFERENCES flavor(id);
ALTER TABLE `flavor_ingredient` ADD FOREIGN KEY (ingredient_id) REFERENCES ingredient(id);
ALTER TABLE `flavor_ingredient` ADD CONSTRAINT uk_flavor_id_ingredient_id UNIQUE (`flavor_id`, `ingredient_id`);

