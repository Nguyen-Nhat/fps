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
	statusInit    = 0
	statusSuccess = 2
	statusFailed  = 3
)