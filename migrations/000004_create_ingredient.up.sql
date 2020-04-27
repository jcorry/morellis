-- Create ingredient table
CREATE TABLE ingredient (
    id INT(11) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;