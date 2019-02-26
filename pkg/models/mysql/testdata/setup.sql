
CREATE TABLE `user` (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
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
    name VARCHAR(16) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `ref_user_status` (id, name)
VALUES
(1, "Unverified"),
(2, "Verified"),
(3, "Deleted");

ALTER TABLE `user` ADD FOREIGN KEY (status_id) REFERENCES ref_user_status(id);

-- Insert dummy user
INSERT INTO user (first_name, last_name, email, phone, status_id, hashed_password, created)
VALUES (
    'Alice',
    'Jones',
    'alice@example.com',
    '867-5309',
    2,
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
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
