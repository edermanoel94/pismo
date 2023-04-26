# Desafio Pismo

## Requirements

- Go 1.19++
- Make
- Docker
- docker-compose

## Building

### Source code

To build the source code, enter the root project folder and execute:

```bash
make build
```

### Run inside a container

```bash
make start-docker
```

## Testing

and to run tests execute
```bash
make test
```

## Setup Environment
Add a `configuration.yml` file in `infra/config` package.

To run it locally, you must add at least the properties below to `configuration.yml` file to execute service:
```yml
server:
  addr: :8080
  debug: false
  timeout:
    read-seconds: 15
    write-seconds: 20

db:
  host: localhost
  user: postgres
  password: pismo
  name: pismo
 
operation_types:
  COMPRA_A_VISTA: "-"
  COMPRA_PARCELADA: "-"
  SAQUE: "-"
  PAGAMENTO: "+"
```

### Folders Structure

This repository contains three main folders: `cmd`, `build` and `internal`.

The `internal` folder contains all the go code, modules and tests that compose the service.

The `build/package` folder contains the Dockerfile used for building the container.

The `cmd/pismo` folder stores the `main.go` file.

The `resources/postman` folder contains all the resources to access each endpoint.

The `resources/sql` folder contains all the sql code that create database.



