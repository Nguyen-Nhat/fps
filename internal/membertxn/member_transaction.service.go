package membertxn

import (
	"context"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	Service interface {
		GetByFileAwardPointId(context.Context, int32) ([]MemberTransaction, error)
		GetByFileAwardPointIDStatuses(context.Context, int32, []int16) ([]MemberTransaction, error)
		Create(context.Context, MemberTxnDTO) (*MemberTransaction, error)
		UpdateOne(context.Context, UpdateMemberTxnDTO) (*MemberTransaction, error)
	}

	ServiceImpl struct {
		repo Repo
	}
)

var _ Service = &ServiceImpl{}

func NewService(repo Repo) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

// Implementation function ---------------------------------------------------------------------------------------------

func (s ServiceImpl) GetByFileAwardPointId(ctx context.Context, fileAwardPointId int32) ([]MemberTransaction, error) {
	return s.repo.FindByFileAwardPointId(ctx, fileAwardPointId)
}
func (s ServiceImpl) GetByFileAwardPointIDStatuses(ctx context.Context, fileAwardPointId int32, statuses []int16) ([]MemberTransaction, error) {
	return s.repo.FindByFileAwardPointIDStatuses(ctx, fileAwardPointId, statuses)
}

func (s ServiceImpl) Create(ctx context.Context, memberTxn MemberTxnDTO) (*MemberTransaction, error) {
	savedMemberTxn, err := s.repo.Save(ctx, MemberTransaction{
		MemberTransaction: ent.MemberTransaction{
			FileAwardPointID: int32(memberTxn.FileAwardPointID),
			Point:            int64(memberTxn.Point),
			Phone:            memberTxn.Phone,
			OrderCode:        memberTxn.OrderCode,
			TxnDesc:          memberTxn.TxnDesc,
			Status:           StatusInit,
			Error:            memberTxn.Error,
			SentTime:         memberTxn.SentTime,
			RefID:            memberTxn.RefID,
			LoyaltyTxnID:     memberTxn.LoyaltyTxnID,
		},
	})
	if err != nil {
		logger.Errorf("Cannot insert file award point, got %v", err)
		return nil, err
	}
	return savedMemberTxn, nil
}
func (s ServiceImpl) UpdateOne(ctx context.Context, memberTxn UpdateMemberTxnDTO) (*MemberTransaction, error) {
	savedMemberTxn, err := s.repo.UpdateOne(ctx, MemberTransaction{
		MemberTransaction: ent.MemberTransaction{
			ID:           int(memberTxn.ID),
			RefID:        memberTxn.RefID,
			SentTime:     memberTxn.SentTime,
			Status:       memberTxn.Status,
			Error:        memberTxn.Error,
			LoyaltyTxnID: memberTxn.LoyaltyTxnID,
		},
	})
	if err != nil {
		logger.Errorf("Cannot insert file award point, got %v", err)
		return nil, err
	}
	return savedMemberTxn, nil
}
