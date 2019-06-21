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
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `flavor` WRITE;
/*!40000 ALTER TABLE `flavor` DISABLE KEYS */;

INSERT INTO `flavor` (`id`, `name`, `description`, `created`, `updated`)
VALUES
	(1,'Coconut Japaleno','One of our most unique flavors, it must be tasted to be believed!\nOur fresh made coconut ice cream is infused with just the right amount of fresh jalapenos. The experience of hot, sweet and cold hits your palate in pretty amazing ways; come try for yourself!','2019-03-01 21:52:22',NULL),
	(2,'Butter Pecan','Butter Pecan is an ice cream standard, but that doesnt mean the flavor has to be ordinary!\nOur buttery, nutty and savory ice cream is a rich and delicious fan favorite, blended with just the right amount of buttery goodness and fresh Georgia pecans.','2019-03-02 21:36:19',NULL);

/*!40000 ALTER TABLE `flavor` ENABLE KEYS */;
UNLOCK TABLES;


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

LOCK TABLES `flavor_store` WRITE;
/*!40000 ALTER TABLE `flavor_store` DISABLE KEYS */;

INSERT INTO `flavor_store` (`id`, `flavor_id`, `store_id`, `position`, `is_active`, `activated`, `deactivated`)
VALUES
	(3,2,1,1,NULL,'2019-03-29 03:08:44','2019-04-07 03:02:54'),
	(4,2,1,1,NULL,'2019-03-29 03:09:00','2019-04-07 03:02:54'),
	(5,2,1,1,NULL,'2019-04-05 13:57:14','2019-04-07 03:02:54'),
	(6,2,1,1,NULL,'2019-04-05 13:57:15','2019-04-07 03:02:54'),
	(7,2,1,1,NULL,'2019-04-05 13:57:17','2019-04-07 03:02:54'),
	(8,2,1,1,NULL,'2019-04-05 13:57:18','2019-04-07 03:02:54'),
	(9,2,1,1,NULL,'2019-04-06 20:48:34','2019-04-07 03:02:54'),
	(12,2,1,1,NULL,'2019-04-06 20:55:10','2019-04-07 03:02:54'),
	(13,2,1,1,NULL,'2019-04-06 20:55:12','2019-04-07 03:02:54'),
	(14,2,1,1,NULL,'2019-04-06 20:55:15','2019-04-07 03:02:54'),
	(15,2,1,1,NULL,'2019-04-07 03:01:49','2019-04-07 03:02:54');

/*!40000 ALTER TABLE `flavor_store` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table ingredient
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ingredient`;

CREATE TABLE `ingredient` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `ingredient` WRITE;
/*!40000 ALTER TABLE `ingredient` DISABLE KEYS */;

INSERT INTO `ingredient` (`id`, `name`, `created`, `updated`)
VALUES
	(1,'coconut','2019-03-01 21:52:22',NULL),
	(2,'jalapeno','2019-03-01 21:52:22',NULL),
	(3,'butter','2019-03-02 21:36:19',NULL),
	(4,'pecan','2019-03-02 21:36:19',NULL),
	(5,'nuts','2019-03-02 21:36:19',NULL);

/*!40000 ALTER TABLE `ingredient` ENABLE KEYS */;
UNLOCK TABLES;


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

LOCK TABLES `permission_user` WRITE;
/*!40000 ALTER TABLE `permission_user` DISABLE KEYS */;

INSERT INTO `permission_user` (`user_id`, `permission_id`, `created`)
VALUES
	(1,7,'2019-04-08 02:36:19'),
    (1,3,'2019-04-08 02:38:43');

/*!40000 ALTER TABLE `permission_user` ENABLE KEYS */;
UNLOCK TABLES;

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
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_email` (`email`),
  KEY `status_id` (`status_id`),
  CONSTRAINT `fk_user_status_id_ref_user_status_id` FOREIGN KEY (`status_id`) REFERENCES `ref_user_status` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `uuid`, `first_name`, `last_name`, `email`, `phone`, `status_id`, `hashed_password`, `created`, `updated`)
VALUES
	(1,'e6fc6b5a-882c-40ba-b860-b11a413ec2df','John','Corry','jcorry@gmail.com','678-592-8804',1,X'243261243132245A5230363157324A6E624877424D66386A535A6E4C75456C623674702F716E4630494C4C5972506C5A46435A536E59642E38365043','2019-03-18 12:24:34',NULL);

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

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
