# Chirpy

## Overview  
Chirpy is a Go API for posting short messages ("chirps"). It includes user authentication, token management, and admin tools.

## Development Technologies  
- **Go** – Backend logic  
- **PostgreSQL** – Database  
- **SQLC** – Generates Go code from SQL  
- **Goose** – Database migrations  

## Project Structure  
```plaintext
internal/app/chirpy  -> Application logic & handlers
internal/database    -> SQL queries & models
internal/auth        -> Authentication utilities
/web                 -> Static assets
main.go              -> App entry point
go.mod               -> Dependencies
sqlc.yaml            -> SQLC configuration
```

## Setup

### Configuring Environment  
Create a .env file with the following variables:  
```sh
DB_URL="postgres://<user>:@localhost:5432/<db-name>?sslmode=disable"
SECRET="<your_jwt_secret>"
PLATFORM="dev"
```

### Running Database Migrations  
```sh
goose -dir internal/database/sql/schema postgres "postgres://user:@localhost:5432/<db-name>" up
```

### Generating Go Code from SQL  
```sh
sqlc generate
```

### Starting the Server  
```sh
go run main.go
```
The server runs at [http://localhost:8080](http://localhost:8080).

## API Endpoints

### User Management

#### Create User  
Create a new user account.  
```sh
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "secret"}'
```

#### Login  
Authenticate a user and obtain tokens.  
```sh
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "secret"}'
```

#### Token Refresh  
Obtain a new access token using a refresh token.  
```sh
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer <refresh_token>"
```

#### Token Revoke  
Revoke an existing token.  
```sh
curl -X POST http://localhost:8080/api/revoke \
  -H "Content-Type: application/json" \
  -d '{"token": "<token_to_revoke>"}'
```

### Chirps

#### Post a Chirp  
Create a new chirp message.  
```sh
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"body": "Hello, world!"}'
```

#### Get All Chirps  
Retrieve all chirps.  
```sh
curl -X GET http://localhost:8080/api/chirps
```

#### Get a Specific Chirp  
Retrieve a single chirp by its ID.  
```sh
curl -X GET http://localhost:8080/api/chirps/<chirpID>
```

#### Delete a Chirp  
Delete a chirp by its ID (requires appropriate authentication).  
```sh
curl -X DELETE http://localhost:8080/api/chirps/<chirpID> \
  -H "Authorization: Bearer <access_token>"
```

### Webhooks

#### Polka Webhook  
Endpoint for handling external webhook events.  
```sh
curl -X POST http://localhost:8080/api/polka/webhooks \
  -H "Content-Type: application/json" \
  -d '{"event": "example_event", "data": {"user_id": <user-uuid>}}'
```

### Admin

#### Get Metrics  
Retrieve application metrics.  
```sh
curl -X GET http://localhost:8080/admin/metrics
```

#### Reset Data  
Reset development data (use with caution).  
```sh
curl -X POST http://localhost:8080/admin/reset
```

### Health Check

#### Check Server Health  
Verify that the API is up and running.  
```sh
curl -X GET http://localhost:8080/api/healthz
```

## License  
This project is licensed under the MIT License.
