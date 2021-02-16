# Dump of table ref_user_status
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ref_user_status`;

CREATE TABLE `ref_user_status` (
    `id` tinyint(3) unsigned NOT NULL,
    `name` varchar(16) NOT NULL,
    `slug` varchar(16) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `ref_user_status` WRITE;
/*!40000 ALTER TABLE `ref_user_status` DISABLE KEYS */;

INSERT INTO `ref_user_status` (`id`, `name`, `slug`)
VALUES
(1,'Unverified','unverified'),
(2,'Verified','verified'),
(3,'Deleted','deleted');

/*!40000 ALTER TABLE `ref_user_status` ENABLE KEYS */;
UNLOCK TABLES;

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `uuid` varchar(36) NOT NULL,
    `first_name` varchar(24) DEFAULT NULL,
    `last_name` varchar(24) DEFAULT NULL,
    `email` varchar(128) DEFAULT NULL,
    `phone` varchar(24) NOT NULL,
    `status_id` tinyint(3) unsigned NOT NULL DEFAULT '1',
    `hashed_password` char(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_phone` (`phone`),
    UNIQUE KEY `uk_user_email` (`email`),
    KEY `status_id` (`status_id`),
    CONSTRAINT `fk_user_status_id_ref_user_status_id` FOREIGN KEY (`status_id`) REFERENCES `ref_user_status` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;