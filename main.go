package main

import (
	"context"
	"flag"
	"log"

	kinde_provider "github.com/AZZO/terraform-provider-kinde/kinde_provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/AZZO/kinde",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), kinde_provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
