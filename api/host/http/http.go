package http

import (
	"github.com/Duke1616/api-demo/api"
	"github.com/Duke1616/api-demo/api/host"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/julienschmidt/httprouter"
)

// Host 模块的 HTTP API服务实例
var Api = &handler{}

type handler struct {
	host host.Service
	log  logger.Logger
}

// Init 初始化时候，以来外部host service的实例
//func (h *handler) Init(srv host.Service) {
//	h.log = zap.L().Named("HOST API")
//	h.host = srv
//}

func (h *handler) Init() error {
	h.log = zap.L().Named("HOST API")

	if api.Host == nil {
		panic("dependence host service is nil")
	}
	h.host = api.Host
	return nil
}

// Registry 把handler实现的方法注册给主路由
func (h *handler) Registry(r *httprouter.Router) {
	r.POST("/hosts", h.CreateHost)
	r.GET("/hosts", h.QueryHost)
}
