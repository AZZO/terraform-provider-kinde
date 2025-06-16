# Kinde Provider

The Kinde provider is used to interact with the Kinde API to manage applications and other resources.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

To use the provider, you need to configure it with your Kinde credentials:

```hcl
terraform {
  required_providers {
    kinde = {
      source = "AZZO/kinde"
    }
  }
}

provider "kinde" {
  issuer_url    = "https://your-domain.kinde.com"
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
}
```

## Resources

### kinde_application

The `kinde_application` resource is used to manage Kinde applications.

#### Example Usage

```hcl
resource "kinde_application" "example" {
  name = "example-application"
  type = "reg"  # Can be "reg", "m2m", or "spa"
}
```

#### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the application.
* `type` - (Required) The type of the application. Must be one of: "reg" (Regular), "m2m" (Machine to Machine), "spa" (Single Page Application).

#### Attributes Reference

The following attributes are exported:

* `application_id` - The unique identifier of the application.
* `client_id` - The client ID of the application.
* `client_secret` - The client secret of the application.

## Development

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To generate or update documentation, run `go generate`.

To build the provider, run `go build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To test the provider, run `go test`. Add the `-v` flag for verbose output.

```shell
$ go test -v ./...
```

To run the full suite of Acceptance tests, run `make testacc`.

```shell
$ make testacc
```

Note: Acceptance tests create real resources, and often cost money to run.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 