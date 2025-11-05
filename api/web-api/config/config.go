package endpoint_config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rapidaai/pkg/configs"
	"github.com/spf13/viper"
)

// Application config structure
type AppConfig struct {
	Name            string                 `mapstructure:"service_name" validate:"required"`
	Version         string                 `mapstructure:"version"`
	Host            string                 `mapstructure:"host" validate:"required"`
	Env             string                 `mapstructure:"env" validate:"required"`
	Secret          string                 `mapstructure:"secret" validate:"required"`
	Port            int                    `mapstructure:"port" validate:"required"`
	LogLevel        string                 `mapstructure:"log_level" validate:"required"`
	PostgresConfig  configs.PostgresConfig `mapstructure:"postgres" validate:"required"`
	RedisConfig     configs.RedisConfig    `mapstructure:"redis" validate:"required"`
	WebhookHost     string                 `mapstructure:"webhook_host" validate:"required"`
	WorkflowHost    string                 `mapstructure:"workflow_host" validate:"required"`
	WebHost         string                 `mapstructure:"web_host" validate:"required"`
	ProviderHost    string                 `mapstructure:"provider_host" validate:"required"`
	IntegrationHost string                 `mapstructure:"integration_host" validate:"required"`
	ExperimentHost  string                 `mapstructure:"experiment_host" validate:"required"`
	EndpointHost    string                 `mapstructure:"endpoint_host" validate:"required"`
	//
	DocumentHost     string                   `mapstructure:"document_host" validate:"required"`
	AssetStoreConfig configs.AssetStoreConfig `mapstructure:"asset_store" validate:"required"`
	OpenSearchConfig configs.OpenSearchConfig `mapstructure:"opensearch" validate:"required"`
}

func (cfg *AppConfig) IsDevelopment() bool {
	return cfg.Env != "production"
}

// reading config and intializing configs for application
func InitConfig() (*viper.Viper, error) {
	vConfig := viper.NewWithOptions(viper.KeyDelimiter("__"))

	vConfig.AddConfigPath(".")
	vConfig.SetConfigName(".env")
	path := os.Getenv("ENV_PATH")
	if path != "" {
		log.Printf("env path %v", path)
		vConfig.SetConfigFile(path)
	}
	vConfig.SetConfigType("env")
	vConfig.AutomaticEnv()
	err := vConfig.ReadInConfig()
	if err != nil {
		log.Printf("Error while reading the config %v", err)
	}

	//
	setDefault(vConfig)
	if err = vConfig.ReadInConfig(); err != nil && !os.IsNotExist(err) {
		log.Printf("Reading from env varaibles.")
	}

	return vConfig, nil
}

func setDefault(v *viper.Viper) {
	// setting all default values
	// keeping watch on https://github.com/spf13/viper/issues/188

}

// Getting application config from viper
func GetApplicationConfig(v *viper.Viper) (*AppConfig, error) {
	var config AppConfig
	err := v.Unmarshal(&config)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	// valdating the app config
	validate := validator.New()
	err = validate.Struct(&config)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return &config, nil
}
