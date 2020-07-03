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

This client serves 2 routes:

- "/": home page, the OAuth2 flow starts from there
- "/callback": this is where Hydra should redirect the user after the authorization flow. At this point, the client should have received an access token.

### Setup Hydra

- It's needed to register the client so that Hydra know it exists and which permission can be given to it.

### Build the provider

TODO change when workflow is better

```
cd provider/src
go mod download
go build -o provider
./provider
```

```
cd example
go mod download
go build -o example
./example
```

### Launch the example

- Launch the hydra containers
- Launch the identity provider
- Launch the client
- Go to http://localhost:3003
- Follow the instructions

### Still needs to be done

- Users management
- Fine-grained permissions (hard + fhir is still not sure about how to do that)
