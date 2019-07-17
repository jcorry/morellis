# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.25)
# Database: morellis
# Generation Time: 2019-04-09 03:09:36 +0000
# ************************************************************

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

# Dump of table flavor
# ------------------------------------------------------------

DROP TABLE IF EXISTS `flavor`;

CREATE TABLE `flavor` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL,
    `description` text,
    `created` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# Dump of table flavor_ingredient
# ------------------------------------------------------------

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

# Dump of table flavor_store
# ------------------------------------------------------------

DROP TABLE IF EXISTS `flavor_store`;

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

# Dump of table ingredient
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ingredient`;

CREATE TABLE `ingredient` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL,
    `created` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# Dump of table permission
# ------------------------------------------------------------

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
(1,'store:read'),
(2,'store:write'),
(3,'user:read'),
(4,'user:write'),
(5,'flavor:read'),
(6,'flavor:write'),
(7,'self:write'),
(8,'self:read');

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

# Dump of table ref_user_status
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ref_user_status`;

LOCK TABLES `ref_user_status` WRITE;

CREATE TABLE `ref_user_status` (
    `id` tinyint(3) unsigned NOT NULL,
    `name` varchar(16) NOT NULL,
    `slug` varchar(16) NOT NULL,
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40000 ALTER TABLE `ref_user_status` DISABLE KEYS */;

INSERT INTO `ref_user_status` (`id`, `name`, `slug`)
VALUES
(1,'Unverified','unverified'),
(2,'Verified','verified'),
(3,'Deleted','deleted');

/*!40000 ALTER TABLE `ref_user_status` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table store
# ------------------------------------------------------------

DROP TABLE IF EXISTS `store`;

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
    `created` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated` DATETIME NULL DEFAULT NULL,
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


# Dump of table user
# ------------------------------------------------------------

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
    `created` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_email` (`email`),
    KEY `status_id` (`status_id`),
    CONSTRAINT `fk_user_status_id_ref_user_status_id` FOREIGN KEY (`status_id`) REFERENCES `ref_user_status` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `ingredient_user` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `ingredient_id` int(11) unsigned NOT NULL,
    `user_id` int(11) unsigned NOT NULL,
    `keyword` varchar(16) DEFAULT NULL,
    `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted` int(8) DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_ingredient_user_ingredient_id_user_id` (`ingredient_id`,`user_id`,`deleted`),
    KEY `fk_ingredient_user_user_id` (`user_id`),
    CONSTRAINT `fk_ingredient_user_ingredient_id` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredient` (`id`),
    CONSTRAINT `fk_ingredient_user_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;