Here is a condensed documentation about Hydra's features you could have to use.

# Endpoints

## Public endpoints

### Get public keys

GET {public}/.well-known/jwks.json

This can be useful to validate JWTs.

### Start auth flow

GET {public}/oauth2/auth

You should preferably use a lib (as golang.org/x/oauth2 for instance) to access this endpoint and launch the authorization flow.

### Get a token

POST {public}/oauth2/token

You should preferably use a lib (as golang.org/x/oauth2 for instance) to access this endpoint and exchange a authorization code for a token.

### Get user info

GET {public}/userinfo

You'll need an authorization header with access token (Bearer {access-token}) to access this route.
The response should contain the information you're authorize to query about the token's user.

## Admin endpoints

### Create a client

POST {admin}/clients

With a body that looks like:

```
{
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

### Update a client

PUT {admin}/clients

With a body that looks like:

```
{
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
