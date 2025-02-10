# OAuth-Tower

Lightweight and simple OAuth 2.0 Authorization Server

## Overview

OAuth-Tower is a lightweight and simple OAuth 2.0 authorization server designed to provide secure and efficient authorization for your applications. It supports various OAuth 2.0 flows and is built using Go and the Echo framework.

## Features

- OAuth 2.0 Authorization Code Flow
- OAuth 2.0 Implicit Flow
- OAuth 2.0 Client Credentials Flow
- JSON Web Token (JWT) support
- JSON Web Key (JWK) management
- Well-known configuration endpoint
- Secure client authentication


## Getting Started

### Prerequisites

- Go 1.23.4 or later

### Installation

1. Clone the repository:

```sh
git clone https://github.com/zercos/oauth-tower.git
cd oauth-tower
```
2. Install dependencies:
```
go mod tidy
```
3. Create a .env file with the necessary environment variables:
```
cp .env.example .env
```
4. Generate RSA keys for JWT signing and place them in the keys directory:
```
openssl genpkey -algorithm RSA -out keys/key.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in keys/key.pem -out keys/key.public.pem
```
Running the Server
To run the server in development mode:
```
./scripts/run_dev.sh
```
The server will start on http://0.0.0.0:8000.

Running Tests
To run the tests:
```
./scripts/tests.sh
```

## API Endpoints
- GET /.well-known/oauth-authorization-server
- GET /oauth/authorization
- POST /oauth/token
- POST /oauth/introspection
- POST /oauth/revocation
- GET /oauth/jwks.json

## License
This project is licensed under the Apache License 2.0 - see the LICENSE file for details.
