CREATE TABLE IF NOT EXISTS `config_mapping`
(
    `id`                   bigint       NOT NULL AUTO_INCREMENT,
    `client_id`            int          NOT NULL,
    `total_tasks`          int          NOT NULL DEFAULT '0',
    `data_start_at_row`    int          NOT NULL DEFAULT '0',
    `require_column_index` varchar(255) NOT NULL,
    `error_column_index`   varchar(255) NOT NULL,
    `created_at`           timestamp    NULL     DEFAULT CURRENT_TIMESTAMP,
    `updated_at`           timestamp    NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`           varchar(255) NOT NULL,
    `data_at_sheet`        varchar(100)          DEFAULT '' COMMENT 'Default is first sheet in file',
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `config_task`
(
    `id`                           bigint       NOT NULL AUTO_INCREMENT,
    `config_mapping_id`            int          NOT NULL,
    `task_index`                   int          NOT NULL,
    `end_point`                    varchar(255) NOT NULL,
    `method`                       varchar(255) NOT NULL,
    `header`                       text         NOT NULL,
    `path_params`                  text         NOT NULL,
    `request_params`               text         NOT NULL,
    `request_body`                 text         NOT NULL,
    `response_success_http_status` int          NOT NULL,
    `response_success_code_schema` varchar(255) NOT NULL,
    `response_message_schema`      varchar(255) NOT NULL,
    `group_by_columns`             varchar(50)           DEFAULT '' COMMENT 'Group by list columns name. Eg: A,B,C',
    `group_by_size_limit`          int                   DEFAULT '0' COMMENT 'Max size of a Group. If exceed, reject file',
    `created_at`                   timestamp    NULL     DEFAULT CURRENT_TIMESTAMP,
    `updated_at`                   timestamp    NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`                   varchar(255) NOT NULL,
    `name`                         varchar(255) NOT NULL DEFAULT 'no name',
    `message_transformations`      text COMMENT 'Format JSON, transform message for displaying in file result',
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `fps_client`
(
    `id`              bigint       NOT NULL AUTO_INCREMENT,
    `client_id`       int          NOT NULL,
    `name`            varchar(255) NOT NULL,
    `description`     varchar(255) NOT NULL,
    `created_at`      timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      timestamp    NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by`      varchar(255) NOT NULL,
    `sample_file_url` varchar(255)      DEFAULT '',
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `processing_file`
(
    `id`                    bigint       NOT NULL AUTO_INCREMENT,
    `client_id`             varchar(255) NOT NULL,
    `display_name`          varchar(255) NOT NULL,
    `file_url`              varchar(255) NOT NULL,
    `result_file_url`       varchar(255)          DEFAULT NULL,
    `status`                smallint     NOT NULL,
    `total_mapping`         int                   DEFAULT '0',
    `stats_total_row`       int                   DEFAULT '0',
    `stats_total_success`   int                   DEFAULT '0',
    `error_display`         text,
    `created_at`            timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`            varchar(255) NOT NULL,
    `updated_at`            timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `file_parameters`       text,
    `stats_total_processed` int                   DEFAULT '0',
    `need_group_row`        tinyint(1)            DEFAULT '0',
    `seller_id`             int                   DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `pf_client_id_IDX` (`client_id`),
    KEY `pf_client_id_status_IDX` (`client_id`, `status`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `processing_file_row`
(
    `id`                bigint    NOT NULL AUTO_INCREMENT,
    `file_id`           bigint    NOT NULL,
    `row_index`         int       NOT NULL,
    `row_data_raw`      text      NOT NULL,
    `task_index`        int       NOT NULL,
    `task_mapping`      text      NOT NULL,
    `task_depends_on`   varchar(20)    DEFAULT NULL,
    `task_request_raw`  text,
    `task_response_raw` text,
    `status`            smallint  NOT NULL,
    `error_display`     text,
    `task_request_curl` text,
    `executed_time`     bigint         DEFAULT NULL,
    `created_at`        timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `group_by_value`    text,
    PRIMARY KEY (`id`),
    KEY `pfr_file_id_IDX` (`file_id`),
    KEY `pfr_file_id_status_IDX` (`file_id`, `status`),
    KEY `pfr_row_index_IDX` (`row_index`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `processing_file_row_group`
(
    `id`                 int          NOT NULL AUTO_INCREMENT,
    `file_id`            int          NOT NULL,
    `task_index`         int          NOT NULL,
    `group_by_value`     text         NOT NULL,
    `total_rows`         int          NOT NULL,
    `row_index_list`     text         NOT NULL,
    `group_request_curl` text         NOT NULL,
    `group_response_raw` text         NOT NULL,
    `status`             smallint     NOT NULL,
    `error_display`      varchar(255) NOT NULL DEFAULT '',
    `executed_time`      int          NOT NULL DEFAULT '0',
    `created_at`         timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`         timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `pfgt_file_id_IDX` (`file_id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_general_ci;
