package main

import (
	"context"
	"flag"
	"log"

	"github.com/AZZO/terraform-provider-kinde/kinde"
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

	err := providerserver.Serve(context.Background(), kinde.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
