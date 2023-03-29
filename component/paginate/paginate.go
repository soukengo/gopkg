package paginate

const (
	PageSizeDefault = 20
	PageSizeMax     = 100
)

// Paginating 分页参数
type Paginating struct {
	PageNo   int32 `json:"page_no"`
	PageSize int32 `json:"page_size"`
}

func (x *Paginating) Limit() int64 {
	if x == nil {
		return PageSizeDefault
	}
	if x.PageSize > 0 && x.PageSize <= PageSizeMax {
		return int64(x.PageSize)
	}
	return PageSizeDefault
}

func (x *Paginating) Offset() int64 {
	if x == nil {
		return 0
	}
	page := x.PageNo
	if page <= 0 {
		page = 1
	}
	return (int64(page - 1)) * x.Limit()
}

// Paginated 分页结果
type Paginated struct {
	Total int64 `json:"total,omitempty"`
}
