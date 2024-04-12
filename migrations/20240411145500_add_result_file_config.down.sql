ALTER TABLE config_mapping DROP COLUMN result_file_config;

ALTER TABLE processing_file_row DROP INDEX pfr_file_id_task_index_idx;
