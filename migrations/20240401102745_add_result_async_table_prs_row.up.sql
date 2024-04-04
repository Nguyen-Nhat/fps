alter table processing_file_row add column result_async text default null comment 'Result async API';
alter table processing_file_row add column start_at timestamp default null comment 'Start time call api';
alter table processing_file_row add column receive_response_at timestamp default null comment 'Receive response at';
alter table processing_file_row add column receive_result_async_at timestamp default null comment 'Receive result async at';
