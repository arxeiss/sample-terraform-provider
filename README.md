# Sample Terraform provider

[Terraform](https://www.terraform.io/) is great for managing infrastructure by the code with config files.
Those might be then saved in git or any other source control management.

Every service, which is manageable by Terraform, must produce its own provider. More about custom providers in [the documentation](https://www.terraform.io/docs/extend/).

In this repo, you can find a dummy server and custom provider communicating with the backend. Easy to start, ready to play.

## How to run

### Prerequisites

- Installed [Go](https://golang.org/) 1.16+ - local server is written in Go, as well as provider for Terraform, as it is now only officially supported language for Terraform SDK.
  - SQLite driver is written in C language so `CGO_ENABLED=1` must be set and supported to start the server (`make start_server` will set this variable for you)
- Installed [Terraform CLI](https://www.terraform.io/downloads.html). Optionally you might install it with [ASDF](https://asdf-vm.com/) - read more about it on [Dev.To](https://dev.to/arxeiss/asdf-manage-multiple-runtime-versions-1fn9) or [Kutac.cz](https://www.kutac.cz/pocitace-a-internety/asdf-rozsiritelny-spravce-multi-verzi-runtime).
- Optionally installed Docker - Local server uses SQLite DB and with `docker-compose` you can start in-browser viewer.

### Running

- `make start_server` to start the HTTP server
  - After the start, there will be printed out `access_token` required by terraform provider
- `make terraform_build` to build Terraform provider
- `cd provider/config` to move to the folder with `*.tf` files
  - `terraform plan` to show plan which will be executed
  - `terraform apply` to execute the plan
  - ... or any other terraform commands
