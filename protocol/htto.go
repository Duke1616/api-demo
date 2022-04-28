package protocol

import (
	"fmt"
	"github.com/Duke1616/api-demo/conf"
	"github.com/infraboard/mcube/http/middleware/cors"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"

	hostAPI "github.com/Duke1616/api-demo/api/host/http"
)

func NewHTTPService() *HTTPService {
	r := httprouter.New()

	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.Addr(),
		Handler:           cors.AllowAll().Handler(r),
	}
	return &HTTPService{
		r:      r,
		l:      zap.L().Named("API"),
		c:      conf.C(),
		server: server,
	}
}

type HTTPService struct {
	r      *httprouter.Router
	l      logger.Logger
	c      *conf.Config
	server *http.Server
}

func (s *HTTPService) Start() error {
	err := hostAPI.Api.Init()
	if err != nil {
		return err
	}

	// 注册
	hostAPI.Api.Registry(s.r)

	// 启动HTTP服务
	s.l.Infof("HTTP服务启动成功, 监听地址: %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service is stopped")
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}

func (s *HTTPService) Stop() error {
	return nil
}
