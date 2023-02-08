CREATE TABLE IF NOT EXISTS `users` (
     `user_id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
     `name` VARCHAR(63),
     `email` VARCHAR(255),
     `password` VARCHAR(255),
     `enabled` BOOLEAN DEFAULT 1,
     `role` ENUM ('user', 'admin') DEFAULT 'user',
     `date_created` TIMESTAMP DEFAULT NOW(),
     `date_updated` TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB;