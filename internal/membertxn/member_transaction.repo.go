package membertxn

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/membertransaction"
	"github.com/mitchellh/mapstructure"
)

type (
	Repo interface {
		FindByFileAwardPointId(context.Context, int32) ([]MemberTransaction, error)
	}

	RepoImpl struct {
		client *ent.Client
	}
)

const dbEngine = "mysql"

var _ Repo = &RepoImpl{} // only for mention that RepoImpl implement Repo

// NewRepo ...
func NewRepo(db *sql.DB) *RepoImpl {
	drv := entsql.OpenDB(dbEngine, db)
	client := ent.NewClient(ent.Driver(drv))
	return &RepoImpl{client: client}
}

// Implementation function ---------------------------------------------------------------------------------------------

func (r RepoImpl) FindByFileAwardPointId(ctx context.Context, fileAwardPointId int32) ([]MemberTransaction, error) {
	txns, err := r.client.MemberTransaction.Query().Where(membertransaction.FileAwardPointID(fileAwardPointId)).All(ctx)
	if err != nil {
		return nil, err
	}

	var res []MemberTransaction
	//if err := mapstructure.Decode(txns, res); err != nil {
	//	return nil, fmt.Errorf("failed decode member transaction list from DB")
	//}

	for _, txn := range txns {
		model := &MemberTransaction{}
		if err := mapstructure.Decode(txn, model); err != nil {
			return nil, fmt.Errorf("failed decode file award point info from DB")
		}
		res = append(res, *model)
	}

	return res, nil
}
