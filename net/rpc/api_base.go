package rpc

const (
	XHeaderLogId = "X-LogId"
)

const (
	ErrNone             = "OK"
	ErrResourceNotFound = "ResourceNotFound"
)

type APIRet struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	// Bytes is filled by the logic with the marshaling result of field Data
	Bytes []byte `json:"-"`
	// LogId is filled by the logic with the response header of X-LogId
	LogId string `json:"-"`
}

type PagerData struct {
	PageNo     int         `json:"pageNo"`
	PageSize   int         `json:"pageSize"`
	TotalCount int64       `json:"totalCount"`
	PageList   interface{} `json:"pageList"`
}

type Pager struct {
	PageNo     int
	PageSize   int
	TotalCount int64
	SortBys    map[string]string
}

func (p *Pager) Pages() int {
	return int((p.TotalCount + int64(p.PageSize) - 1) / int64(p.PageSize))
}
