# Identity Provider

This repository is a identity provider designed to work with [Hydra](https://www.ory.sh/hydra/docs/).
Having these makes it possible to distribute OAuth2 authorization grants.

## Routes exposed

- "/login": the user should give its credentials here, it's how we make sure that we'll eventually give access to the API to an authorized actor.
- "/consent": on this page, the user is supposed to choose which scopes it will let the client use when fetching the resource provider.

## OAuth2 flows

### OAuth2 Authorize Code Flow

https://www.ory.sh/hydra/docs/implementing-consent/

Our usecase being to authenticate users and grant us access or not to our API, this is the one that we are the most interested in.

### Other flows

Other flows as trading an access token for a refresh token are possible but still not tested nor showcased in the example.

## Toy showcase

Here, we show how to run Hydra, a identity provider and a basic client on your machine.

The client serves 2 routes:

- "/": home page, the OAuth2 flow starts from there
- "/callback": this is where Hydra should redirect the user after the authorization flow. At this point, the client should have received an access token.

### Setup Hydra

- First, we need to launch the Hydra containers:

```
cd hydra
docker-compose -f compose-hydra.yml -f compose-postgres.yml up --build
```

- Then, we need to register the client so that Hydra knows it exists and which permission can be given to it. We can do so via Hydra's API:

```
POST /clients
Content-Type: application/json
Accept: application/json
Body: {
    "client_id": "open-id-client",
    "client_secret": "secret",
    "grant_types": [
        "authorization_code", "refresh_token"
    ],
    "response_types": [
        "code", "id_token"
    ],
    "scope": "openid offline_access",
    "redirect_uris": [
        "http://localhost:3003/callback"
    ]
}
```

### Build and run the provider

The provider expects a postgresql database to run (credentials for this DB are to be provided in .env).

This DB should contain the users' credentials in 3 columns: email, name and password.

For convenience, the provider exposes a `/signup` route you can POST to to register a new user. To do so, you have to POST a body with fields `name`, `password` and `email`.

Building and running the provider can be done with

```
make run-provider
```

### Build and run the client

```
make run-client
```

### Run through the example

- Go to http://localhost:3003
- Follow the instructions
