alter table config_mapping
    add max_file_size int not null default 5 comment 'Max file size (MB) that client can upload',
    add tenant_id varchar(255) null comment 'Tenant Id of client. Eg: OMNI, CDP, CARBON, ...',
    add using_merchant_attr_name tinyint(1)
        not null default 0 comment 'If 1, when import/get data, FPS will filter by sellerId, platformId,... (based on merchant_attribute_name value)',
    add merchant_attribute_name varchar(255) null comment 'Attribute name of users attribute in IAM that is used for filtering data',
    add ui_config text null comment 'UI config for client. Eg: show hide elements, change positions, ... Ref: https://confluence.teko.vn/display/SupplyChain/%5BFPS%5D+UI+Config';

alter table fps_client
    add import_file_template_url varchar(255) null comment 'Private URL of template file that client can download, need send token when get data in FE';
