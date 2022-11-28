package pagination

type PaginatingRequest struct {
	Page int
	Size int
}

func (s *PaginatingRequest) InitDefaultValue() {
	if s.Page == 0 {
		s.Page = 1
	}
	if s.Size == 0 {
		s.Size = 10
	}
}
