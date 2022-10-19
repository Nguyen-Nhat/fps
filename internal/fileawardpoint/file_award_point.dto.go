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

type CreateFileAwardPointReqDTO struct {
	MerchantID  int64
	FileUrl     string
	Note        string
	FileName    string
	CreatedUser string
}

type CreateFileAwardPointResDTO struct {
	FileAwardPointId int32
}
