package dto

type PaginationRequest struct {
	PageSize  int    `json:"pageSize" form:"pageSize"`
	Page      int    `json:"p" form:"p"`
	Field     string `json:"field" form:"field"`
	Value     string `json:"value" form:"value"`
	SortField string `json:"sortField" form:"sortField"`
	SortOrder string `json:"sortOrder" form:"sortOrder"`
	Owner     string `json:"owner" form:"owner"`
	GroupName string `json:"groupName" form:"groupName"`
}

func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationRequest) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

type PaginationResponse struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type Response struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Data2  interface{} `json:"data2"`
}

func SuccessResponse(data interface{}) *Response {
	return &Response{
		Status: "ok",
		Msg:    "",
		Data:   data,
	}
}

func SuccessResponseWithTotal(data interface{}, total int64) *Response {
	return &Response{
		Status: "ok",
		Msg:    "",
		Data:   data,
		Data2:  total,
	}
}

func ErrorResponse(msg string) *Response {
	return &Response{
		Status: "error",
		Msg:    msg,
		Data:   nil,
	}
}
