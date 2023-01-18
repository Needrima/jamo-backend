package helper

import (
	"github.com/spf13/viper"
	"log"
)

type ConfigStruct struct {
	ServiceAddress     string  `mapstructure:"service_address"`
	ServicePort        string  `mapstructure:"service_port"`
	ServiceMode        string  `mapstructure:"service_mode"`
	DbType             string  `mapstructure:"db_type"`
	MongoDbHost        string  `mapstructure:"mongo_db_host"`
	MongoDbName        string  `mapstructure:"mongo_db_name"`     // database name
	MongoDbUserName    string  `mapstructure:"mongo_db_username"` // authentication username
	MongoDbPassword    string  `mapstructure:"mongo_db_password"` // authentication password
	MongoDbConnString  string  `mapstructure:"mongo_db_conn_string"`
	MongoDbAuthDb      string  `mapstructure:"mongo_db_auth_db"`
	ServiceName        string  `mapstructure:"service_name"`
	LogFile            string  `mapstructure:"log_file"`
	LogDir             string  `mapstructure:"log_dir"`
	ExternalConfigPath string  `mapstructure:"external_config_path"`
	PageLimit          string  `mapstructure:"page_limit"`
	TotalExpense       float64 `mapstructure:"total_expense"`
	SMTPHost           string  `mapstructure:"smtp_host"`
	SMTPPort           string  `mapstructure:"smtp_port"`
	SMTPUsername       string  `mapstructure:"smtp_username"`
	SMTPPassword       string  `mapstructure:"smtp_password"`
}

func LoadEnv(path string) (config ConfigStruct, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("jamo-backend")
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

var Config = ReturnConfig()
