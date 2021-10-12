provider "sdc" {
  endpoint = var.endpoint
  api_token = var.api_token
}

resource "sdc_hdd" "abc" {
    name = "abc"
    description = "jkl"
}
