alter table config_mapping
    drop column max_file_size,
    drop column tenant_id,
    drop column using_merchant_attr_name,
    drop column merchant_attribute_name,
    drop column ui_config;

alter table fps_client
    drop column import_file_template_url;
