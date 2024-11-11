ALTER TABLE processing_file
    MODIFY COLUMN accept_language varchar(5) DEFAULT 'en' NULL COMMENT 'Language of user when upload file (detected by Accept-Language header). Eg: en, vi, ...';

UPDATE processing_file
SET accept_language = 'vi'
WHERE accept_language IS NULL OR accept_language = '';