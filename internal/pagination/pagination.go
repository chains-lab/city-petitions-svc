package pagination

type Request struct {
	Page uint64 `json:"page"`
	Size uint64 `json:"size"`
}

type Response struct {
	Page  uint64 `json:"page"`
	Size  uint64 `json:"size"`
	Total uint64 `json:"total"`
}

func CalculateLimitOffset(req Request) (limit uint64, offset uint64) {
	limit = req.Size
	offset = (req.Page - 1) * req.Size

	if limit == 0 {
		limit = 10 // default limit if not specified
	}
	if offset < 0 {
		offset = 0 // ensure offset is not negative
	}

	return
}
