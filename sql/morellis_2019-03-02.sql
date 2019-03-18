# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.21)
# Database: morellis
# Generation Time: 2019-03-02 21:57:38 +0000
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
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `description` text,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated` timestamp NULL DEFAULT NULL,
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
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `flavor_id` int(11) NOT NULL,
  `ingredient_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_flavor_id_ingredient_id` (`flavor_id`,`ingredient_id`),
  KEY `ingredient_id` (`ingredient_id`),
  CONSTRAINT `flavor_ingredient_ibfk_1` FOREIGN KEY (`flavor_id`) REFERENCES `flavor` (`id`),
  CONSTRAINT `flavor_ingredient_ibfk_2` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredient` (`id`)
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


# Dump of table ingredient
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ingredient`;

CREATE TABLE `ingredient` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
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
	(1,'Unverified', 'unverified'),
	(2,'Verified', 'verified'),
	(3,'Deleted', 'deleted');

/*!40000 ALTER TABLE `ref_user_status` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table store
# ------------------------------------------------------------

DROP TABLE IF EXISTS `store`;

CREATE TABLE `store` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
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



# Dump of table user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uuid` VARCHAR(36) NOT NULL,
  `first_name` varchar(24) DEFAULT NULL,
  `last_name` varchar(24) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  `phone` varchar(24) NOT NULL,
  `status_id` tinyint(3) unsigned NOT NULL DEFAULT '1',
  `hashed_password` CHAR(60) NOT NULL,d
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_email` (`email`),
  KEY `status_id` (`status_id`),
  CONSTRAINT `user_ibfk_1` FOREIGN KEY (`status_id`) REFERENCES `ref_user_status` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `first_name`, `last_name`, `email`, `phone`, `status_id`, `hashed_password`, `created`, `updated`)
VALUES
	(1,'John','Corry','jcorry@gmail.com','678-592-8804',1,'$2a$12$mDwXnc5gFLEr3gOZTi7Keuenwhtzs9xPQ.RLfg9HkKSlrqEJgmka.','2019-03-01 18:12:25',NULL);

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
