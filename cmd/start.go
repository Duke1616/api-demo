package cmd

import (
	"errors"
	"fmt"
	"github.com/Duke1616/api-demo/api"
	"github.com/Duke1616/api-demo/api/host/impl"
	"github.com/Duke1616/api-demo/conf"
	"github.com/Duke1616/api-demo/protocol"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	configType string
	configFile string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Demo后端API服务",
	Long:  `Demo后端API服务`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载全局配置
		if err := loadGlobalConfig(configType); err != nil {
			return err
		}

		// 初始化日志
		if err := loadGlobalLogger(); err != nil {
			return err
		}

		// 初始化服务层
		if err := impl.Service.Init(); err != nil {
			return err
		}

		// 把服务实例注册给IOC层
		api.Host = impl.Service

		// 启动服务后需要处理的事件
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

		// 启动服务
		srv := NewService(conf.C())

		// 等待程序退出
		go srv.waitSign(ch)

		// 启动服务
		if err := srv.Start(); err != nil {
			if !strings.Contains(err.Error(), "http: Server closed") {
				return err
			}
		}
		return nil
	},
}

func NewService(conf *conf.Config) *Service {
	return &Service{
		conf: conf,
		http: protocol.NewHTTPService(),
	}
}

type Service struct {
	conf *conf.Config
	http *protocol.HTTPService
}

func (s *Service) Start() error {
	return s.http.Start()
}

// 当发现用户手动终止，需要处理
func (s *Service) waitSign(ch chan os.Signal) {

}

func loadGlobalConfig(configType string) error {
	switch configType {
	case "file":
		err := conf.LoadConfigFromToml(configFile)
		if err != nil {
			return err
		}
	case "env":
		err := conf.LoadConfigFromEnv()
		if err != nil {
			return err
		}
	case "etcd":
		return errors.New("not implemented")
	default:
		return errors.New("unknown config type")
	}
	return nil
}

func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)
	lc := conf.C().Log
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}
	zapConfig := zap.DefaultConfig()
	zapConfig.Level = level
	zapConfig.Files.RotateOnStartup = false
	switch lc.To {
	case conf.ToStdout:
		zapConfig.ToStderr = true
		zapConfig.ToFiles = false
	case conf.ToFile:
		zapConfig.Files.Name = "api.log"
		zapConfig.Files.Path = lc.PathDir
	}
	switch lc.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}

func init() {
	startCmd.PersistentFlags().StringVarP(&configFile, "config_file", "f", "etc/restful-api.toml", "the restful-api config file path")
	startCmd.PersistentFlags().StringVarP(&configType, "config_type", "t", "file", "the restful-api config type")
	rootCmd.AddCommand(startCmd)
}
