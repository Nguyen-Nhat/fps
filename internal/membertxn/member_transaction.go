package membertxn

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"time"
)

// MemberTransaction ... is model of table `member_transaction`
type MemberTransaction struct {
	ent.MemberTransaction
}

func Name() string {
	return "MemberTransaction"
}

func (t *MemberTransaction) IsCheckExpires(expiresTimeMinutes int) bool {
	now := time.Now()
	expireAt := t.CreatedAt.Local().Add(time.Minute * time.Duration(expiresTimeMinutes))
	return now.After(expireAt)
}

// Status ENUM ...
const (
	StatusInit       = 1
	StatusProcessing = 2
	StatusFailed     = 3
	StatusSuccess    = 4
	StatusTimeout    = 5
)
