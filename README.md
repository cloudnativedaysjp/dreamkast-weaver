# dreamkast-weaver

Backend Server for dreamkast/dreamkast-ui.

<div align="center">
<img src="./images/icon.jpg" alt="dreamkast-weaver" width="300">
</div>

## Prerequisites

- Docker
- Docker Compose
- [Service Weaver](https://serviceweaver.dev/docs.html#installation)

## How to run

Run the dev container and database:

```bash
docker-compose -f dev/compose.yaml up -d
```

Access `http://localhost:8080` then you can perform tests with a graphiql UI.

Dev app container supports live-reloading since it is running on [air](https://github.com/cosmtrek/air).
You don't need to rebuild the container image except for the case `go.mod` updated.
