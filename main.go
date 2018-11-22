package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/zircote/terraform-provider-aurora/aurora"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aurora.Provider})
}
