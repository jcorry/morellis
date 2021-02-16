DROP TABLE IF EXISTS `flavor_ingredient`;

CREATE TABLE `flavor_ingredient` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `flavor_id` int(11) unsigned NOT NULL,
    `ingredient_id` int(11) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_flavor_id_ingredient_id` (`flavor_id`,`ingredient_id`),
    KEY `ingredient_id` (`ingredient_id`),
    CONSTRAINT `fk_flavor_ingredient_fid_flavor_id` FOREIGN KEY (`flavor_id`) REFERENCES `flavor` (`id`),
    CONSTRAINT `fk_flavor_ingredient_iid_ingredient_id` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredient` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `flavor_ingredient` WRITE;
/*!40000 ALTER TABLE `flavor_ingredient` DISABLE KEYS */;

INSERT INTO `flavor_ingredient` (`id`, `flavor_id`, `ingredient_id`)
VALUES
(1,1,1),
(2,1,2),
(3,2,3),
(4,2,4),
(5,2,5);

/*!40000 ALTER TABLE `flavor_ingredient` ENABLE KEYS */;
UNLOCK TABLES;