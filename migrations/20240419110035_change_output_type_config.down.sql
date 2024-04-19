alter table config_mapping
    modify output_file_type enum ('XLSX', 'CSV') default 'XLSX' not null comment 'Type of file output ' after input_file_type;
