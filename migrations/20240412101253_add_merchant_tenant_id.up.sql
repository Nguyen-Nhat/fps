alter table processing_file add column tenant_id varchar(20) comment 'Identifying the tenant of the request. Its value can be CDP, OMNI, etc';
alter table processing_file add column merchant_id varchar(20) comment 'Identifying the merchant in the tenant. platformId (from CDP) or sellerId (from OMNI)';
update processing_file set merchant_id = cast(seller_id as char(20)) where 1;
