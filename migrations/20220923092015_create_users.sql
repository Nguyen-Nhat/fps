-- create "users" table
CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT, `name` varchar(255) NOT NULL DEFAULT 'unknown', `active` bool NOT NULL DEFAULT true, `email` varchar(255) NOT NULL, `phone_number` varchar(255) NOT NULL, `password_hash` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `user_id` (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
