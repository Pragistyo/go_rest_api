### GO- PROJECT


## STACK

Language: GO . \
Server/router: (gorilla/mux). \
Database: PostgreSQL (with pgx)\

## Feature:


jsonwebtoken: authentication(done) and authorization (backlog)

## ROUTES

#### User Routes

| Route                                |  HTTP  | Description |
| ------------------------------------ | ------ | --------------|
| `/api/v1/user`        | GET    | Get all user data
| `/api/v1/user`        | POST   | Create one user data
| `/api/v1/user/{id}/`  | GET    | Get one user data
| `/api/v1/user/{id}/`  | PUT    | Update one user data
| `/api/v1/user/{id}/`  | DELETE | Delete one user data
| `/api/v1/login/`      | POST   | Login user, return token
| `/api/v1/verified/`  | GET    | Retrieve user data based on token
