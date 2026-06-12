package main

import (
	"github.com/martinweber/steampipe-plugin-pve/pve"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: pve.Plugin})
}
