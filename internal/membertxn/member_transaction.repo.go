package membertxn

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/membertransaction"
)

type (
	Repo interface {
		FindByFileAwardPointId(context.Context, int32) ([]MemberTransaction, error)
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