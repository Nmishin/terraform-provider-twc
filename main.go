package main

import (
    "context"
    "log"

    "terraform-provider-twc/internal/provider"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
    opts := providerserver.ServeOpts{
        Address: "registry.terraform.io/Nmishin/twc",
    }

    err := providerserver.Serve(context.Background(), provider.New(), opts)
    if err != nil {
        log.Fatal(err.Error())
    }
}
