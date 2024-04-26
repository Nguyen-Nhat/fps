alter table config_mapping
    modify output_file_type enum ('XLS', 'XLSX', 'CSV') null comment 'Type of file output (XLS, XLSX, CSV). If null, output type is input type. If has value will force output type' after input_file_type;
