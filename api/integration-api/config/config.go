package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rapidaai/pkg/configs"
	"github.com/spf13/viper"
)

// Application config structure
type AppConfig struct {
	Name           string                 `mapstructure:"service_name" validate:"required"`
	Version        string                 `mapstructure:"version" validate:"required"`
	Secret         string                 `mapstructure:"secret" validate:"required"`
	Host           string                 `mapstructure:"host" validate:"required"`
	Port           int                    `mapstructure:"port" validate:"required"`
	LogLevel       string                 `mapstructure:"log_level" validate:"required"`
	PostgresConfig configs.PostgresConfig `mapstructure:"postgres" validate:"required"`
	RedisConfig    configs.RedisConfig    `mapstructure:"redis" validate:"required"`
	SendgridApiKey string                 `mapstructure:"sendgrid_api_key" validate:"required"`

	AssetStoreConfig configs.AssetStoreConfig `mapstructure:"asset_store" validate:"required"`

	// all the host which is required
	ProviderHost    string `mapstructure:"provider_host" validate:"required"`
	WebHost         string `mapstructure:"web_host" validate:"required"`
	IntegrationHost string `mapstructure:"integration_host" validate:"required"`
	EndpointHost    string `mapstructure:"endpoint_host" validate:"required"`
	ExperimentHost  string `mapstructure:"experiment_host" validate:"required"`
	WebhookHost     string `mapstructure:"webhook_host" validate:"required"`
	WorkflowHost    string `mapstructure:"workflow_host" validate:"required"`
	DocumentHost    string `mapstructure:"document_host" validate:"required"`
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
	if err == nil {
		log.Printf("Error while reading the config")
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

	v.SetDefault("SERVICE_NAME", "go-service-template")
	v.SetDefault("VERSION", "0.0.1")
	v.SetDefault("HOST", "0.0.0.0")
	v.SetDefault("PORT", 9090)
	v.SetDefault("LOG_LEVEL", "debug")
	v.SetDefault("SENDGRID_API_KEY", "")

	// all internal service host
	v.SetDefault("PROVIDER_HOST", "")
	v.SetDefault("VAULT_HOST", "")
	v.SetDefault("INTEGRATION_HOST", "")
	v.SetDefault("PROJECT_HOST", "")
	//

	v.SetDefault("POSTGRES__HOST", "localhost")
	v.SetDefault("POSTGRES__PORT", 5432)
	v.SetDefault("POSTGRES__DB_NAME", "<>")
	v.SetDefault("POSTGRES__AUTH__USER", "<>")
	v.SetDefault("POSTGRES__AUTH__PASSWORD", "<>")
	v.SetDefault("POSTGRES__MAX_OPEN_CONNECTION", 10)
	v.SetDefault("POSTGRES__MAX_IDEAL_CONNECTION", 10)
	v.SetDefault("POSTGRES__SSL_MODE", "disable")
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
