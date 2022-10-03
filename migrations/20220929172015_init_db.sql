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
) CHARSET UTF8MB4
  COLLATE UTF8MB4_BIN;
