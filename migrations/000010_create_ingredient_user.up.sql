-- Creat ingredient_user table
CREATE TABLE ingredient_user (
    id int(11) unsigned NOT NULL AUTO_INCREMENT,
    ingredient_id int(11) unsigned NOT NULL,
    user_id int(11) unsigned NOT NULL,
    keyword varchar(16) DEFAULT NULL,
    created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted int(8) DEFAULT '0',
    PRIMARY KEY (id),
    UNIQUE KEY uk_ingredient_user_ingredient_id_user_id (ingredient_id,user_id,deleted),
    KEY fk_ingredient_user_user_id (user_id),
    CONSTRAINT fk_ingredient_user_ingredient_id FOREIGN KEY (ingredient_id) REFERENCES ingredient (id),
    CONSTRAINT fk_ingredient_user_user_id FOREIGN KEY (user_id) REFERENCES user (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;