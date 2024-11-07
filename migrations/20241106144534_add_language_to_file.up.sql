ALTER TABLE processing_file
    ADD accept_language varchar(3) NULL DEFAULT '' COMMENT 'Language of user when upload file (detected by Accept-Language header). Eg: en, vi, ...';
