terraform {
  required_providers {
    kinde = {
      source = "CSymes/kinde"
    }
  }
}

provider "kinde" {
    issuer_url = var.issuer_url
    client_id = var.client_id
    client_secret = var.client_secret
}



data "kinde_application" "test" {
    application_id = var.reference_id
}

resource "kinde_application" "test2" {
    name = "Foobar"
    type = "reg"
}



output "test" {
    value = data.kinde_application.test
}
