-- Create flavor_ingredient table
CREATE TABLE flavor_ingredient (
    id INT(11) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    flavor_id INT(11) UNSIGNED NOT NULL,
    FOREIGN KEY (flavor_id) REFERENCES flavor(id),
    ingredient_id INT(11) UNSIGNED  NOT NULL,
    FOREIGN KEY (ingredient_id) REFERENCES ingredient(id),
    CONSTRAINT uk_flavor_id_ingredient_id UNIQUE (flavor_id, ingredient_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;