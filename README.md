# dreamkast-weaver

Backend Server for dreamkast/dreamkast-ui.

<div align="center">
<img src="./images/icon.jpg" alt="dreamkast-weaver" width="300">
</div>

## Prerequisites

- Docker

## How to run

Run the dev container and database:

```bash
docker compose up -d
```

Access `http://localhost:8088` then you can perform tests with a graphiql UI.

For example, you can calculate the number of viewers using the graphql scripts below:

the viewing script

```graphql
mutation {
  viewTrack(input: {
    profileID: 123,
    trackName: "A",
    talkID: 456
  })
}
```

and the retrieval script.

```graphql
query {
  viewerCount(confName: cndt2023) {
    trackName
    count
  }
}
```
