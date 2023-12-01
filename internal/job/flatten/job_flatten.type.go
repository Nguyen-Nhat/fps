package flatten

type ErrorRow struct {
	RowId  int // id file, start at 0
	Reason string
	//RowData []string
}

const SellerIDKey = "sellerId"
