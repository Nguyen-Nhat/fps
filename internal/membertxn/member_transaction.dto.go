package membertxn

import "time"

type MemberTxnDTO struct {
	ID               int64
	FileAwardPointID int64
	Point            int64
	Phone            string
	OrderCode        string
	RefID            string
	SentTime         time.Time
	TxnDesc          string
	Status           int16
	Error            string
	LoyaltyTxnID     int64
}

type UpdateMemberTxnDTO struct {
	ID           int64
	RefID        string
	SentTime     time.Time
	Status       int16
	Error        string
	LoyaltyTxnID int64
}
