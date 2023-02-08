CREATE TABLE sales (
    `sales_id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
    `user_id`      BIGINT unsigned,
    `product_id`   BIGINT unsigned,
    `quantity`     INT,
    `paid`         INT,
    `date_created` TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (sales_id),
    CONSTRAINT sale_user_fk FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL,
    CONSTRAINT sale_product_fk FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE SET NULL
 );