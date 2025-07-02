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
	Status    int    `json:"status"`
	Code      string `json:"code"`
	Message   string `json:"msg"`
	Data      any    `json:"data"`
	RequestID string `json:"requestID"`
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
	Page    int   `json:"page"`
	PerPage int   `json:"perPage"`
	Total   int64 `json:"total"`
	Items   []any `json:"items"`
}

type Pager struct {
	Page    int
	PerPage int
	Total   int64
	SortBys map[string]string
}

func (p *Pager) Pages() int {
	return int((p.Total + int64(p.PerPage) - 1) / int64(p.Page))
}

func (p *Pager) Set(page, perPage int, total int64) {
	p.Page = page
	p.PerPage = perPage
	p.Total = total
}
