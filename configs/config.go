package configs

import "github.com/spf13/viper"

type conf struct {
	RedisHost        string `mapstructure:"REDIS_HOST"`
	RedisPort        string `mapstructure:"REDIS_PORT"`
	RedisDb          int    `mapstructure:"REDIS_DB"`
	RateLimitDefault int    `mapstructure:"RATE_LIMIT_DEFAULT"`
	TimeBlockDefault int    `mapstructure:"TIME_BLOCK_DEFAULT"`
	WebServerPort    string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.SetConfigFile(path + "/.env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
