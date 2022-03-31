package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Environment  string `envconfig:"ENV"` // prod,dev
	Host         string `envconfig:"HOST" default:"https://boosters.news.com"`
	HTTPHostAddr string `envconfig:"HTTP_HOST_ADDR" default:"127.0.0.1:8095"`
	SecretKey    string `envconfig:"GDPR_APP_SECRET"` // <hex>
	JWTSecretKey string `envconfig:"JWT_SECRET"`      // <hex>

	Databases
	S3storage
	Emails
}

type S3storage struct {
	AccessKey string `envconfig:"S3_SPACES_ACCESS_KEY"` // token
	SecretKey string `envconfig:"S3_SPACES_SECRET_KEY"` // token
	Bucket    string `envconfig:"S3_SPACES_BUCKET"`     // costless-prod
	Region    string `envconfig:"S3_SPACES_REGION" default:"fra1"`
	Endpoint  string `envconfig:"S3_SPACES_ENDPOINT" default:"digitaloceanspaces.com"`
}

type Databases struct {
	PostgresURL  string `envconfig:"POSTGRES_URL"`
	ReIndexerURL string `envconfig:"REINDEXER_URL"` // <reindexer_host_url>
}

type Emails struct {
	Admin string `envconfig:"ADMIN_EMAIL"  default:"admin@boosters.news.com"`
}

func NewConfig() *Config {
	cfg := Config{}
	err := envconfig.Process("BOOSTERS", &cfg)
	if err != nil {
		panic(err.Error())
	}
	return &cfg
}

func (c *Config) Debug() bool {
	return c.Environment == "debug"
}
