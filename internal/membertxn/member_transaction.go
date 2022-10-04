package membertxn

import "time"

// MemberTransaction ... is model of table `member_transaction`
type MemberTransaction struct {
	Id               int       `json:"id"`
	FileAwardPointId int       `json:"file_award_point_id"`
	Point            int       `json:"point"`
	Phone            string    `json:"phone"`
	OrderCode        string    `json:"order_code"`
	RefId            string    `json:"ref_id"`
	SentTime         time.Time `json:"sent_time"`
	TxnDesc          string    `json:"txn_desc"`
	Status           int       `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Status ENUM ...
const (
	statusInit    = 0
	statusSuccess = 2
	statusFailed  = 3
)
