
# walletsvc

Service to handle Wallet & Transaction

## Features

- Init Wallet Account
- Enable Wallet
- Create Deposit
- Create Withdraw
- Disable Wallet
- Get Wallet
- Get Wallet Transaction




## Tech Stack


- Golang 1.17
- Redis 6.2
- MySQL 8.0
- Docker


## Run Locally Docker

Clone the project

```bash
  git clone https://github.com/alvinatthariq/walletsvc
```

Go to the project directory

```bash
  cd walletsvc
```

Install dependencies

- Docker https://docs.docker.com/desktop/install/mac-install/

Run docker compose

```bash
  docker-compose up -d
```

Server will be running at localhost:8080


## Run Locally

Clone the project

```bash
  git clone https://github.com/alvinatthariq/walletsvc
```

Go to the project directory

```bash
  cd walletsvc
```

Install dependencies

- MySQL 8.0
- Redis 6.2




Start the server

```bash
  make run
```

Server will be running at localhost:8080


## Running Tests

To run tests, run the following command

```bash
  make run-test
```

## Documentation

[Documentation](https://documenter.getpostman.com/view/27910682/2s93z9agqK)