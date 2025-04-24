package config

import (
  "github.com/spf13/viper"
)

type RemoteCustomOption struct {
	Contain string `mapstructure:"contain"`
	Options RemoteOptions `mapstructure:"options"`
}

type RemoteOptions struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	InsecureTLS bool `mapstructure:"insecureTLS"`
}

func (r *RemoteOptions) Merge (r1 RemoteOptions){
	if(r1.Username != ""){
		r.Username = r1.Username
	}
	if(r1.Password != ""){
		r.Password = r1.Password
	}
	if(r1.InsecureTLS){
		r.InsecureTLS=r1.InsecureTLS
	}
}

type Config struct {
	RemoteCustomOption []RemoteCustomOption `mapstructure:"remoteCustomOption"`

}

func GetConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") 
	viper.AddConfigPath(".")    
	viper.SetDefault("RemoteCustomOption",[]RemoteCustomOption{})
	viper.ReadInConfig()
	var config Config
	viper.Unmarshal(&config)
	return config
}
