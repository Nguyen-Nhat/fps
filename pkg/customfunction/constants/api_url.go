package constants

const (
	// api-doc: https://apidoc.teko.vn/project-doc/approved/core_logic_layer/file_service_retail/version/latest/operations/post_uploads
	UrlFileServiceUploadImage = "http://files-core-api.files-service/upload/image" // for calling internal service -> should move to env config

	// api-doc: https://apidoc.teko.vn/project-doc/approved/core_logic_layer/retail/3_order_processing/warehouse_central/version/latest/operations/GetSites
	UrlApiGetSites = "http://warehouse3-central-service.warehouse-management/api/v1/sites"

	// api-doc: https://apidoc.teko.vn/project-doc/approved/core_logic_layer/retail/catalog/version/latest/paths/api-v2-skus/get
	UrlApiGetSkus = "http://catalog-core-v2-api.catalog/api/v2/skus" // url call service name
)
