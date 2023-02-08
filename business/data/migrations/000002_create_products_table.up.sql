CREATE TABLE IF NOT EXISTS products (
 `product_id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `name`         TEXT,
  `cost`        INT,
  `quantity`     INT,
  `user_id`     BIGINT unsigned,
 `date_created` TIMESTAMP DEFAULT NOW(),
 `date_updated` TIMESTAMP DEFAULT NOW(),

  PRIMARY KEY (product_id),
  CONSTRAINT product_user_fk FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL
);