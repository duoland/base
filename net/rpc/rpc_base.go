package rpc

const (
	XHeaderLogID = "X-Request-ID"
)

const (
	ErrNone             = "OK"
	ErrResourceNotFound = "ResourceNotFound"
)

type APIRet interface {
	IsOk() bool
	HasResp() bool
	RetCode() string
	RetMessage() string
	SetRequestID(string)
}

type BaseAPIRet struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	RequestID string `json:"requestId"`
}

func (r *BaseAPIRet) SetRequestID(requestID string) {
	r.RequestID = requestID
}

func (r *BaseAPIRet) IsOk() bool {
	return r.Code == ErrNone
}

func (r *BaseAPIRet) HasResp() bool {
	return r.Code != ""
}

func (r *BaseAPIRet) RetCode() string {
	return r.Code
}

func (r *BaseAPIRet) RetMessage() string {
	return r.Message
}

type PagerData struct {
	PageNo     int   `json:"pageNo"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	PageList   any   `json:"pageList"`
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

func (p *Pager) Set(pageNo, pageSize int, totalCount int64) {
	p.PageNo = pageNo
	p.PageSize = pageSize
	p.TotalCount = totalCount
}
