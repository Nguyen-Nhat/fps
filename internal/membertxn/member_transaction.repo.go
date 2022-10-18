package membertxn

import (
	"context"
	"database/sql"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/membertransaction"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Repo interface {
		FindByFileAwardPointId(context.Context, int32) ([]MemberTransaction, error)
		FindByFileAwardPointIDStatuses(context.Context, int32, []int16) ([]MemberTransaction, error)
		Save(context.Context, MemberTransaction) (*MemberTransaction, error)
		UpdateOne(context.Context, MemberTransaction) (*MemberTransaction, error)
	}

	repoImpl struct {
		client *ent.Client
	}
)

const dbEngine = "mysql"

var _ Repo = &repoImpl{} // only for mention that repoImpl implement Repo

// NewRepo ...
func NewRepo(db *sql.DB) *repoImpl {
	drv := entsql.OpenDB(dbEngine, db)
	client := ent.NewClient(ent.Driver(drv))
	return &repoImpl{client: client}
}

// Implementation function ---------------------------------------------------------------------------------------------

func (r repoImpl) FindByFileAwardPointId(ctx context.Context, fileAwardPointId int32) ([]MemberTransaction, error) {
	txnArr, err := r.client.MemberTransaction.Query().
		Where(membertransaction.FileAwardPointID(fileAwardPointId)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var res []MemberTransaction
	for _, txn := range txnArr {
		model := &MemberTransaction{*txn}
		res = append(res, *model)
	}

	return res, nil
}
func (r repoImpl) FindByFileAwardPointIDStatuses(ctx context.Context, fileAwardPointId int32, statuses []int16) ([]MemberTransaction, error) {
	txnArr, err := r.client.MemberTransaction.Query().
		Where(membertransaction.FileAwardPointID(fileAwardPointId)).
		Where(membertransaction.StatusIn(statuses...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var res []MemberTransaction
	for _, txn := range txnArr {
		model := &MemberTransaction{*txn}
		res = append(res, *model)
	}

	return res, nil
}

func (r *repoImpl) Save(ctx context.Context, memberTxn MemberTransaction) (*MemberTransaction, error) {
	return save(ctx, r.client, memberTxn)
}
func (r *repoImpl) UpdateOne(ctx context.Context, memberTxn MemberTransaction) (*MemberTransaction, error) {
	return update(ctx, r.client, memberTxn)
}

func save(ctx context.Context, client *ent.Client, memberTxn MemberTransaction) (*MemberTransaction, error) {
	// 1. Create
	memberTxnSaved, err := mapMemberTxn(client, memberTxn).Save(ctx)
	if err != nil {
		logger.Errorf("Cannot save to member transaction, got: %v", err)
		return nil, fmt.Errorf("failed to save member transaction to DB")
	}

	// 2. Return
	return &MemberTransaction{*memberTxnSaved}, nil
}
func update(ctx context.Context, client *ent.Client, memberTxn MemberTransaction) (*MemberTransaction, error) {
	// 1. Update
	memberTxnSaved, err := client.MemberTransaction.
		UpdateOneID(memberTxn.ID).
		SetRefID(memberTxn.RefID).
		SetSentTime(memberTxn.SentTime).
		SetLoyaltyTxnID(memberTxn.LoyaltyTxnID).
		SetTxnDesc(memberTxn.TxnDesc).
		SetStatus(memberTxn.Status).
		SetError(memberTxn.Error).
		Save(ctx)
	if err != nil {
		logger.Errorf("Cannot update to member transaction :%#v, got: %v", memberTxn, err)
		return nil, fmt.Errorf("failed to update member transaction to DB")
	}

	// 2. Return
	return &MemberTransaction{*memberTxnSaved}, nil
}

// mapMemberTxn ... must update this function when schema is modified
func mapMemberTxn(client *ent.Client, memberTxn MemberTransaction) *ent.MemberTransactionCreate {
	return client.MemberTransaction.Create().
		SetFileAwardPointID(memberTxn.FileAwardPointID).
		SetPoint(memberTxn.Point).
		SetPhone(memberTxn.Phone).
		SetOrderCode(memberTxn.OrderCode).
		SetRefID(memberTxn.RefID).
		SetSentTime(memberTxn.SentTime).
		SetTxnDesc(memberTxn.TxnDesc).
		SetStatus(memberTxn.Status).
		SetError(memberTxn.Error)
}
