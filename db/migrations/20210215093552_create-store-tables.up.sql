CREATE TABLE `store` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL,
    `phone` varchar(32) DEFAULT NULL,
    `email` varchar(32) DEFAULT NULL,
    `url` varchar(64) DEFAULT NULL,
    `address` varchar(128) DEFAULT NULL,
    `city` varchar(64) DEFAULT NULL,
    `state` varchar(32) DEFAULT NULL,
    `zip` varchar(16) DEFAULT NULL,
    `lat` decimal(9,6) DEFAULT NULL,
    `lng` decimal(9,6) DEFAULT NULL,
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `store` WRITE;
/*!40000 ALTER TABLE `store` DISABLE KEYS */;

INSERT INTO `store` (`id`, `name`, `phone`, `email`, `url`, `address`, `city`, `state`, `zip`, `lat`, `lng`, `created`, `updated`)
VALUES
(1,'Morellis On Moreland','404-622-0210','info@morellisicecream.com','http://www.morellisicecream.com/','749 Moreland Ave SE','Atlanta','GA','30316',33.733951,-84.349625,'2019-03-27 00:15:25',NULL),
(2,'Dunwoody Farmburger','404-622-0210','info@morellisicecream.com','http://www.morellisicecream.com/','4514 Chamblee Dunwoody Rd','Dunwoody','GA','30338',33.922714,-84.315169,'2019-03-27 00:31:17',NULL);

/*!40000 ALTER TABLE `store` ENABLE KEYS */;
UNLOCK TABLES;

CREATE TABLE `flavor_store` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `flavor_id` int(11) unsigned NOT NULL,
    `store_id` int(11) unsigned NOT NULL,
    `position` smallint(6) NOT NULL,
    `is_active` tinyint(1) DEFAULT '0',
    `activated` datetime NOT NULL,
    `deactivated` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_flavor_store_is_active_store_id_position_id` (`store_id`,`position`,`is_active`),
    KEY `fk_flavor_store_flavor_id_flavor_id` (`flavor_id`),
    CONSTRAINT `fk_flavor_store_flavor_id_flavor_id` FOREIGN KEY (`flavor_id`) REFERENCES `flavor` (`id`),
    CONSTRAINT `fk_flavor_store_store_id_store_id` FOREIGN KEY (`store_id`) REFERENCES `store` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `flavor_store` WRITE;
/*!40000 ALTER TABLE `flavor_store` DISABLE KEYS */;