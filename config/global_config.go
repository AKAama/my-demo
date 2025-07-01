package config

import (
	"github.com/go-viper/mapstructure/v2"
	"myapi/pkg/util"
	"os"
	"path/filepath"
	_ "regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type IConfig interface {
	Validate() []error
}

type GlobalConfig struct {
	Port     int       `json:"port,omitempty" yaml:"port,omitempty"`
	DBConfig *DBConfig `json:"db" yaml:"db"`
}

func (g *GlobalConfig) Validate() []error {
	var errs = make([]error, 0)
	if err := util.IsValidPort(g.Port); err != nil {
		errs = append(errs, err)
	}
	if es := g.DBConfig.Validate(); len(es) > 0 {
		errs = append(errs, es...)
	}
	return errs
}

func NewDefaultGlobalConfig() *GlobalConfig {
	cfg := &GlobalConfig{
		Port:     3000,
		DBConfig: NewDefaultDBConfig(),
	}
	return cfg
}

func TryLoadFromDisk(configFilePath string) (*GlobalConfig, error) {
	_, err := os.Stat(configFilePath)
	if err != nil {
		return nil, err
	}
	dir, file := filepath.Split(configFilePath)
	fileType := filepath.Ext(file)
	viper.AddConfigPath(dir)
	viper.SetConfigName(strings.TrimSuffix(file, fileType))
	viper.SetConfigType(strings.TrimPrefix(fileType, "."))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, err
		}
		return nil, errors.Errorf("解析配置文件错误:%s", err.Error())
	}
	cfg := NewDefaultGlobalConfig()
	if err := viper.Unmarshal(cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = strings.TrimPrefix(fileType, ".")
	}); err != nil {
		return nil, err
	}
	return cfg, nil
}
