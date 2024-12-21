CREATE DATABASE IF NOT EXISTS gorder;
USE gorder;

DROP TABLE IF EXISTS `o_stock`;

CREATE TABLE `o_stock` (
                           id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                           product_id VARCHAR(255) NOT NULL,
                           quantity INT UNSIGNED NOT NULL DEFAULT 0,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO o_stock (product_id, quantity)
VALUES ('prod_RK6REfs9s1WHHx', 1000);

