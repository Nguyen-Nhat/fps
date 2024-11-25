update config_task
set request_body = '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"supplierId","type":"integer","valuePattern":"$func.convertSupplierCode2SupplierId;$param.sellerId;$A","required":true},{"field":"siteIds","type":"string","valuePattern":"$func.convertSiteCodes2SiteIds;$param.sellerId;$C","required":true},{"field":"isConsignment","type":"boolean","valuePattern":"$func.convertString2Bool;$D","required":true},{"field":"fromDate","type":"string","valuePattern":"$func.convertDateTimeFormat;$E;02/01/2006;2006-01-02T15:04:05Z","required":true},{"field":"toDate","type":"string","valuePattern":"$func.convertDateTimeFormat;$F;02/01/2006;2006-01-02T15:04:05Z","required":true},{"field":"createdSource","type":"string","valuePattern":"import","required":true},{"field":"requestedById","type":"string","valuePattern":"$param.requestedById","required":true},{"field":"action","type":"string","valuePattern":"UPDATE_END_DATE","required":true},{"field":"lineItems","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$func.convertSellerSkuAndUomName2Sku;$param.sellerId;$H;$J","required":false},{"field":"sellerSku","type":"string","valuePattern":"$H","required":true},{"field":"startDate","type":"string","valuePattern":"$func.convertDateTimeFormat;$E;02/01/2006;2006-01-02T15:04:05Z","required":true},{"field":"endDate","type":"string","valuePattern":"$func.convertDateTimeFormat;$F;02/01/2006;2006-01-02T15:04:05Z","required":true},{"field":"price","type":"number","valuePattern":"$H","required":true},{"field":"priceAfterTax","type":"number","valuePattern":"$I","required":true}]}]'
where config_mapping_id = 29
  and task_index = 1;

update config_mapping
set error_column_index = '$J',
    data_at_sheet      = 'data_converted',
    data_start_at_row  = 2
where id = 29;
