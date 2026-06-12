package pve

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type pveConfig struct {
	ApiUrl         *string `hcl:"api_url"`
	User           *string `hcl:"user"`
	Password       *string `hcl:"password"`
	ApiTokenID     *string `hcl:"api_token_id"`
	ApiTokenSecret *string `hcl:"api_token_secret"`
	TLSInsecure    *bool   `hcl:"tls_insecure"`
}

func ConfigInstance() interface{} {
	return &pveConfig{}
}

func GetConfig(connection *plugin.Connection) pveConfig {
	if connection == nil || connection.Config == nil {
		return pveConfig{}
	}
	config, _ := connection.Config.(pveConfig)
	return config
}
