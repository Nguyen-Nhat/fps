package membertxn

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

// MemberTransaction ... is model of table `member_transaction`
type MemberTransaction struct {
	ent.MemberTransaction
}

// Status ENUM ...
const (
	StatusInit       = 1
	StatusProcessing = 2
	StatusFailed     = 3
	StatusSuccess    = 4
)
