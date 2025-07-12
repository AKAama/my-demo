package cmd

import (
	"context"
	"errors"

	"myapi/config"
	"myapi/pkg/db"
	"myapi/pkg/server"
	"myapi/pkg/signals"
	"myapi/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func NewRootCommand() *cobra.Command {
	var configFilePath string
	cmd := &cobra.Command{
		Use:   "",
		Short: "",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableNoDescFlag:   true,
			DisableDescriptions: true,
			HiddenDefaultCmd:    true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.TryLoadFromDisk(configFilePath)
			if err != nil {
				zap.S().Errorf("读取本地配置文件错误:%s", err.Error())
				return
			}
			if errs := cfg.Validate(); len(errs) > 0 {
				zap.S().Errorf("本地配置文件验证错误:%s", errors.Join(errs...))
				return
			}
			ctx := signals.SetupSignalHandler()
			if err := db.InitTiDB(cfg); err != nil {
				zap.S().Infof("数据库连接错误:%s", err.Error())
				return
			}
			zap.S().Infof("数据库连接成功，地址为：%s:%d,库名为：%s", cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.Database)
			if err := run(cfg, ctx); err != nil {
				zap.S().Errorf("运行时错误:%s", err.Error())
				return
			}
		},
		Version: util.GetVersion().Version,
	}
	cmd.Flags().StringVarP(&configFilePath, "config", "c", "./etc/config.yaml", "配置文件路径")
	_ = cmd.MarkFlagRequired("config")
	_ = viper.BindPFlag("config", cmd.Flags().Lookup("config"))
	return cmd
}

func run(cfg *config.GlobalConfig, ctx context.Context) error {
	s := server.NewServer(cfg)
	g, c := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.Run()
	})
	zap.S().Debugf("http server[:%d] 已经运行...", cfg.Port)
	g.Go(func() error {
		<-c.Done()
		s.GracefulShutdown(ctx)
		return nil
	})
	return g.Wait()

}
