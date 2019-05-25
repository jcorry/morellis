
CREATE TABLE `user` (
    id INT(11) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    uuid VARCHAR(36) NOT NULL,
    first_name VARCHAR(24) NULL,
    last_name VARCHAR(24) NULL,
    email VARCHAR(128) NULL,
    phone VARCHAR(24) NOT NULL,
    status_id TINYINT(3) UNSIGNED NOT NULL DEFAULT 1,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `user` ADD CONSTRAINT uk_user_email UNIQUE(email);

CREATE TABLE `ref_user_status` (
    id TINYINT(3) UNSIGNED NOT NULL PRIMARY KEY,
    name VARCHAR(16) NOT NULL,
    slug VARCHAR(16) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `ref_user_status` (id, name, slug)
VALUES
(1, "Unverified", "unverified"),
(2, "Verified", "verified"),
(3, "Deleted", "deleted");

ALTER TABLE `user` ADD FOREIGN KEY (status_id) REFERENCES ref_user_status(id);

-- Insert dummy user
INSERT INTO user (uuid, first_name, last_name, email, phone, status_id, hashed_password, created)
VALUES (
   UUID(),
   'Alice',
   'Jones',
   'alice@example.com',
   '867-5309',
   2,
   '$2a$12$3/ZmSDnMwcfRcxgahkkNrOzyOv28HXbEu1vSVdkIbFId.AKJryrDC', -- 'password'
   '2019-02-24 17:25:25'
);

-- Create syntax for TABLE 'permission'
CREATE TABLE `permission` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4;

INSERT INTO `permission` (`id`, `name`)
VALUES
(1, 'store:read'),
(2, 'store:write'),
(3, 'user:read'),
(4, 'user:write'),
(5, 'flavor:read'),
(6, 'flavor:write'),
(7, 'self:write'),
(8, 'self:read');

-- Create syntax for TABLE 'permission_user'
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

INSERT INTO `permission_user` (`permission_id`, `user_id`)
VALUES
(3, 1),
(4, 1);

-- CREATE store table
CREATE TABLE `store` (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    phone VARCHAR(32) NULL,
    email VARCHAR(32) NULL,
    url VARCHAR(64) NULL,
    address VARCHAR(128) NULL,
    city VARCHAR(64) NULL,
    state VARCHAR(32) NULL,
    zip VARCHAR(16) NULL,
    lat DECIMAL(9, 6),
    lng DECIMAL(9, 6),
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `store` (`id`, `name`, `phone`, `email`, `url`, `address`, `city`, `state`, `zip`, `lat`, `lng`, `created`, `updated`) VALUES (1, 'Morellis On Moreland', '404-622-0210', 'info@morellisicecream.com', 'http://www.morellisicecream.com/', '749 Moreland Ave SE', 'Atlanta', 'GA', '30316', 33.733951, -84.349625, '2019-03-27 00:15:25', NULL);
INSERT INTO `store` (`id`, `name`, `phone`, `email`, `url`, `address`, `city`, `state`, `zip`, `lat`, `lng`, `created`, `updated`) VALUES (2, 'Dunwoody Farmburger', '404-622-0210', 'info@morellisicecream.com', 'http://www.morellisicecream.com/', '4514 Chamblee Dunwoody Rd', 'Dunwoody', 'GA', '30338', 33.922714, -84.315169, '2019-03-27 00:31:17', NULL);

-- Create ingredient table
CREATE TABLE `ingredient` (
    id INT(11) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `ingredient` (`id`, `name`, `created`, `updated`)
VALUES
(1, 'coconut', '2019-03-01 21:52:22', NULL),
(2, 'jalapeno', '2019-03-01 21:52:22', NULL),
(3, 'butter', '2019-03-02 21:36:19', NULL),
(4, 'pecan', '2019-03-02 21:36:19', NULL),
(5, 'nuts', '2019-03-02 21:36:19', NULL);

-- Creat ingredient_user table
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

-- Create flavor table
CREATE TABLE `flavor` (
    id INT(11) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    description TEXT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP null
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `flavor` (`id`, `name`, `description`, `created`, `updated`)
VALUES
(1, 'Coconut Japaleno', 'One of our most unique flavors, it must be tasted to be believed!\nOur fresh made coconut ice cream is infused with just the right amount of fresh jalapenos. The experience of hot, sweet and cold hits your palate in pretty amazing ways; come try for yourself!', '2019-03-01 21:52:22', NULL),
(2, 'Butter Pecan', 'Butter Pecan is an ice cream standard, but that doesnt mean the flavor has to be ordinary!\nOur buttery, nutty and savory ice cream is a rich and delicious fan favorite, blended with just the right amount of buttery goodness and fresh Georgia pecans.', '2019-03-02 21:36:19', NULL);


-- Create flavor_ingredient table
CREATE TABLE `flavor_ingredient` (
    id INT(11) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    flavor_id INT(11) UNSIGNED NOT NULL,
    ingredient_id INT(11) UNSIGNED  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `flavor_ingredient` ADD FOREIGN KEY (flavor_id) REFERENCES flavor(id);
ALTER TABLE `flavor_ingredient` ADD FOREIGN KEY (ingredient_id) REFERENCES ingredient(id);
ALTER TABLE `flavor_ingredient` ADD CONSTRAINT uk_flavor_id_ingredient_id UNIQUE (`flavor_id`, `ingredient_id`);

INSERT INTO `flavor_ingredient` (`id`, `flavor_id`, `ingredient_id`)
VALUES
(1, 1, 1),
(2, 1, 2),
(3, 2, 3),
(4, 2, 4),
(5, 2, 5);

CREATE TABLE `flavor_store` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `flavor_id` int(11) NOT NULL,
    `store_id` int(11) NOT NULL,
    `position` smallint(6) NOT NULL,
    `is_active` tinyint(1) NULL DEFAULT '0',
    `activated` datetime NOT NULL,
    `deactivated` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_flavor_store_is_active_store_id_position_id` (`store_id`,`position`,`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;