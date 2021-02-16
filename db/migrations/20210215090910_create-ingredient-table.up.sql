CREATE TABLE `ingredient` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL,
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `ingredient` WRITE;
/*!40000 ALTER TABLE `ingredient` DISABLE KEYS */;

INSERT INTO `ingredient` (`id`, `name`)
VALUES
(1,'coconut'),
(2,'jalapeno'),
(3,'butter'),
(4,'pecan'),
(5,'nuts');

/*!40000 ALTER TABLE `ingredient` ENABLE KEYS */;
UNLOCK TABLES;
