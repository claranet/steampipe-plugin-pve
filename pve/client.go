package pve

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/luthermonson/go-proxmox"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*proxmox.Client, error) {
	cacheKey := "pve_client"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*proxmox.Client), nil
	}

	config := GetConfig(d.Connection)

	if config.ApiUrl == nil {
		return nil, fmt.Errorf("api_url must be configured in the connection")
	}

	opts := make([]proxmox.Option, 0)

	insecure := config.TLSInsecure != nil && *config.TLSInsecure
	if insecure {
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		opts = append(opts, proxmox.WithHTTPClient(httpClient))
	}

	// Token-based authentication (preferred)
	if config.ApiTokenID != nil && config.ApiTokenSecret != nil {
		opts = append(opts, proxmox.WithAPIToken(*config.ApiTokenID, *config.ApiTokenSecret))
		client := proxmox.NewClient(*config.ApiUrl, opts...)
		d.ConnectionManager.Cache.Set(cacheKey, client)
		return client, nil
	}

	// Username/password authentication
	if config.User != nil && config.Password != nil {
		client := proxmox.NewClient(*config.ApiUrl, opts...)
		_, err := client.Ticket(ctx, &proxmox.Credentials{
			Username: *config.User,
			Password: *config.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to authenticate with username/password: %w", err)
		}
		d.ConnectionManager.Cache.Set(cacheKey, client)
		return client, nil
	}

	return nil, fmt.Errorf("either api_token_id + api_token_secret or user + password must be configured")
}
