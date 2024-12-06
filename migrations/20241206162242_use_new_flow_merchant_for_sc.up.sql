update config_mapping set tenant_id = 'OMNI', using_merchant_attr_name = 1, merchant_attribute_name = 'seller_id' where id in (10,11,12,16,20,27,28,29,33);
update processing_file set tenant_id = 'OMNI', merchant_id = seller_id where client_id in (10,11,12,16,20,27,28,29,33);
