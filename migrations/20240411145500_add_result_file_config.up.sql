ALTER TABLE config_mapping ADD COLUMN result_file_config varchar(500) DEFAULT '' COMMENT 'JSON string, config new column in file result for display process result';

CREATE INDEX pfr_file_id_task_index_idx USING BTREE ON processing_file_row (file_id,task_index);
