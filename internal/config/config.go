package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var Config struct {
	Progstack ProgstackParams `mapstructure:"progstack"`
	Github    GithubParams    `mapstructure:"github"`
	Db        DbParams        `mapstructure:"postgres"`
}

type ProgstackParams struct {
	LoginRedirectUri string `mapstructure:"login_redirect_uri"`
}

type GithubParams struct {
	AppName       string `mapstructure:"app_name"`
	ClientID      string `mapstructure:"client_id"`
	ClientSecret  string `mapstructure:"client_secret"`
	WebhookSecret string `mapstructure:"webhook_secret"`
}

type DbParams struct {
	Name     string `mapstructure:"name"`
	Schema   string `mapstructure:"schema"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
}

func (params DbParams) Connect() (*sql.DB, error) {
	connstr := fmt.Sprintf(
		"user=%s password=%s port=%d dbname=%s sslmode=disable",
		params.User,
		params.Password,
		params.Port,
		params.Name,
	)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
