CREATE TABLE flavor_store (
    id int(11) unsigned NOT NULL AUTO_INCREMENT,
    flavor_id int(11) NOT NULL,
    store_id int(11) NOT NULL,
    position smallint(6) NOT NULL,
    is_active tinyint(1) NULL DEFAULT '0',
    activated datetime NOT NULL,
    deactivated datetime DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_flavor_store_is_active_store_id_position_id (store_id,position,is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;