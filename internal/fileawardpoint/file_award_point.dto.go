package fileawardpoint

// GetFileAwardPointDetailDTO ...
type GetFileAwardPointDetailDTO struct {
	Id int
}

type GetListFileAwardPointDTO struct {
	MerchantId int
	Page       int
	Size       int
}

func (s *GetListFileAwardPointDTO) InitDefaultValue() {
	if s.Page == 0 {
		s.Page = 1
	}
	if s.Size == 0 {
		s.Size = 10
	}
}
