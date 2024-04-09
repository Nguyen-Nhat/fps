alter table config_mapping
    add input_file_type varchar(255) default 'XLSX' not null comment 'Các định dạng cho phép của file input, cách nhau bằng dấu phẩy (ex: "XLSX,CSV") ' after error_column_index;

alter table config_mapping
    add output_file_type enum ('XLSX', 'CSV') default 'XLSX' not null comment 'Type of file output ' after input_file_type;
