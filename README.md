# Chat

Chat application using Go, Websocket, Graphql, Clean Architecture

## Techstack

- Go
- [gqlgen](https://github.com/99designs/gqlgen) for Graphql in Go
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

# Install

You have to install redis in local or using docker to run this

- `make install-tools` to install tools
- `make serve` to run the app

Currently there not yet Frontend side for this app yet.
You can check by using Graphql tools like https://www.postman.com/ or https://insomnia.rest/
Or just browse directly to localhost:8080/ and play with the playground

# Setup Enviroment Variable

| name                 | meaning                                                                         |
| -------------------- | ------------------------------------------------------------------------------- |
| FIREBASE_CREDENTIALS | Firebase credentials, can get when register Firebase and convert to json string |
| DATASOURCE_HOST      | database host                                                                   |
| DATASOURCE_USER      | database user                                                                   |
| DATASOURCE_PASS      | database chat                                                                   |
| DATASOURCE_PORT      | database pass                                                                   |
| DATASOURCE_DATABASE  | database name                                                                   |
| DEBUG_USER_ID        | firebase debug user id                                                          |
| CORS_ALLOWED_ORIGINS | allowed cors host                                                               |
| PORT                 | port                                                                            |

# References

- https://outcrawl.com/go-graphql-realtime-chat

# TODO

- Create Frontend
