-- Create flavor table
CREATE TABLE flavor (
    id INT(11) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    description TEXT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;