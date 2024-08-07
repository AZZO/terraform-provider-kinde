package main

import (
	"github.com/CSymes/terraform-provider-kinde/kinde"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kinde.Provider,
	})
}
