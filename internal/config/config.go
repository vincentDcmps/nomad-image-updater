package config

import (
	"github.com/spf13/viper"
)

type Git struct {
	Enabled   bool   `mapstructure:"enabled"`
	RefBranch string `mapstructure:"refbranch"`
}
type LoggerOption struct {
	Level string `mapstructure:"level"`
}
type GetTagReplaceURL struct {
	Target  string `mapstructure:"target"`
	Replace string `mapstructure:"replace"`
}
type RemoteCustomOption struct {
	Contain string        `mapstructure:"contain"`
	Options RemoteOptions `mapstructure:"options"`
}

type RemoteOptions struct {
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	InsecureTLS bool   `mapstructure:"insecureTLS"`
}

func (r *RemoteOptions) Merge(r1 RemoteOptions) {
	if r1.Username != "" {
		r.Username = r1.Username
	}
	if r1.Password != "" {
		r.Password = r1.Password
	}
	if r1.InsecureTLS {
		r.InsecureTLS = r1.InsecureTLS
	}
}

type Config struct {
	GetTagReplaceURL   []GetTagReplaceURL
	RemoteCustomOption []RemoteCustomOption `mapstructure:"remoteCustomOption"`
	LoggerOption       LoggerOption
	Git                Git
}

var configPath = []string{
	".",
	"~/.config/",
	"/etc/",
}

func GetConfig() Config {
	viper.SetConfigName("nomad-image-updater")
	viper.SetConfigType("yaml")
	for _, v := range configPath {
		viper.AddConfigPath(v)
	}
	viper.SetDefault("RemoteCustomOption", []RemoteCustomOption{})
	viper.SetDefault("LoggerOption.Level", "DEBUG")
	viper.SetDefault("Git.enabled", false)
	viper.SetDefault("Git.refbranch", "master")
	viper.ReadInConfig()
	viper.SetEnvPrefix("NID")
	viper.AutomaticEnv()
	var config Config
	viper.Unmarshal(&config)
	return config
}
