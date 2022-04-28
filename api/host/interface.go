package host

import "context"

type Service interface {
	// CreateHost 录入主机信息
	CreateHost(context.Context, *Host) error
	// QueryHost 查看主机列表信息
	QueryHost(context.Context, *QueryHostRequest) (*Set, error)
	// DescribeHost 主机详情查询
	DescribeHost(context.Context, *DescribeHostRequest) (*Host, error)
	// UpdateHost 主机信息修改
	UpdateHost(context.Context, *UpdateHostRequest) (*Host, error)
	// DeleteHost 删除主机
	DeleteHost(context.Context, *DeleteHostRequest) (*Host, error)
}

type QueryHostRequest struct {
	PageSize   int
	PageNumber int
}

func (req *QueryHostRequest) Offset() int {
	return (req.PageNumber - 1) * req.PageSize
}

type DescribeHostRequest struct {
	Id string
}

type UpdateHostRequest struct {
	Id string
}

type DeleteHostRequest struct {
	Id string
}
