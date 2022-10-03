-- create "file_award_point" table
DROP TABLE IF EXISTS `file_award_point`;

CREATE TABLE `file_award_point`
(
    `id`                  INT          NOT NULL AUTO_INCREMENT,
    `merchant_id`         BIGINT       NOT NULL,
    `display_name`        VARCHAR(255) NOT NULL,
    `file_url`            VARCHAR(255) NOT NULL,
    `result_file_url`     VARCHAR(255)          DEFAULT NULL,
    `status`              INT          NOT NULL DEFAULT 0,
    `stats_total_row`     INT          NOT NULL DEFAULT 0,
    `stats_total_success` INT          NOT NULL DEFAULT 0,
    `created_at`          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`          VARCHAR(45)           DEFAULT NULL,
    `updated_by`          VARCHAR(45)           DEFAULT NULL,
    PRIMARY KEY (`id`)
);

-- create "member_transaction" table
DROP TABLE IF EXISTS `member_transaction`;

CREATE TABLE `member_transaction`
(
    `id`                  INT         NOT NULL AUTO_INCREMENT,
    `file_award_point_id` INT         NOT NULL,
    `point`               INT         NOT NULL,
    `phone`               VARCHAR(15) NOT NULL,
    `order_code`          VARCHAR(50) NOT NULL,
    `ref_id`              VARCHAR(50)          DEFAULT NULL,
    `sent_time`           TIMESTAMP            DEFAULT NULL,
    `txn_desc`            VARCHAR(255)         DEFAULT NULL,
    `status`              INT         NOT NULL DEFAULT 0,
    `created_at`          TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          TIMESTAMP            DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`          VARCHAR(45)          DEFAULT NULL,
    `updated_by`          VARCHAR(45)          DEFAULT NULL,
    PRIMARY KEY (`id`)
);
