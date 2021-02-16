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