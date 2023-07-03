package helper

import (
	"log"

	"github.com/spf13/viper"
)

type ConfigStruct struct {
	ServiceAddress     string `mapstructure:"Service__Address"`
	ServicePort        string `mapstructure:"Service__Port"`
	ServiceMode        string `mapstructure:"Service__Mode"`
	ServiceName        string `mapstructure:"Service__Name"`
	LogFile            string `mapstructure:"Service__LogFileName"`
	LogDir             string `mapstructure:"Service__LogDirectory"`
	LaunchUrl          string `mapstructure:"Service__LaunchUrl"`
	AppName            string `mapstructure:"Service__AppName"`
	Account            string `mapstructure:"Service__Account"`
	Key                string `mapstructure:"Service__Key"`
	DBConnectionString string `mapstructure:"DBConnection__ConnectionString"`
	DBName             string `mapstructure:"DBConnection__DatabaseName"`
	PageLimit          string `mapstructure:"DBConnection__PageLimit"`
	DBConnectionType   string `mapstructure:"DBConnection__Type"`
	EBConnectionString string `mapstructure:"EBConnection__ConnectionString"`
	EBConnectionTTL    string `mapstructure:"EBConnection__TTl"`
	ExternalConfigPath string `mapstructure:"external_config_path"`
	NotificationExpiry string `mapstructure:"Service__NotificationExpiry"`
	GoogleSmtpUsername string `mapstructure:"Google__SmtpUsername"`
	GoogleSmtpPassword string `mapstructure:"Google__SmtpPassword"`
	GoogleSmtpHost     string `mapstructure:"Google__SmtpHost"`
	GoogleSmtpPort     string `mapstructure:"Google__SmtpPort"`
}

func LoadEnv(path string) (config ConfigStruct, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("walls-notification-service")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {

		return ConfigStruct{}, err
	}
	err = viper.Unmarshal(&config)
	return
}
func ReturnConfig() ConfigStruct {
	config, err := LoadEnv(".")
	if err != nil {
		log.Println(err)
	}
	if config.ExternalConfigPath != "" {
		viper.Reset()
		config, err = LoadEnv(config.ExternalConfigPath)
		if err != nil {

			log.Println(err)
		}
	}
	return config
}

var ServiceConfiguration = ReturnConfig()
