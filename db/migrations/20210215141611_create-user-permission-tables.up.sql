DROP TABLE IF EXISTS `permission`;

CREATE TABLE `permission` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `permission` WRITE;
/*!40000 ALTER TABLE `permission` DISABLE KEYS */;

INSERT INTO `permission` (`id`, `name`)
VALUES
(1, 'store:read'),
(2, 'store:write'),
(3, 'user:read'),
(4, 'user:write'),
(5, 'flavor:read'),
(6, 'flavor:write'),
(7, 'ingredient:read'),
(8, 'ingredient:write'),
(9, 'self:write'),
(10, 'self:read'),
(11, 'all');

/*!40000 ALTER TABLE `permission` ENABLE KEYS */;
UNLOCK TABLES;

# Dump of table permission_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `permission_user`;

CREATE TABLE `permission_user` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` int(11) unsigned NOT NULL,
    `permission_id` int(11) unsigned NOT NULL,
    `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_permission_user_permission_id_user_id` (`user_id`,`permission_id`),
    KEY `idx_permission_user_user_id` (`user_id`),
    CONSTRAINT `fk_user_permission_permission_id_permission_id` FOREIGN KEY (`id`) REFERENCES `permission` (`id`),
    CONSTRAINT `fk_user_permission_user_id_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;