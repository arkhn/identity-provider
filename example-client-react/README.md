# Example of a React client app

To run this example locally, here are the steps you have to follow:

- ## 1. Launch the oauth provider (Hydra)

    To run the oauth provider, execute the following in the folder `hydra`.

    ```
    docker-compose -f compose-hydra.yml -f compose-postgres.yml up
    ```

- ## 2. Register the client

    You can register the example client with the following query to Hydra:

    ```json
    POST http://localhost:4445/clients
    {
        "client_id": "example-client-react",
        "client_secret": "example-client-react",
        "grant_types": [
            "authorization_code", "refresh_token"
        ],
        "response_types": [
            "code", "id_token"
        ],
        "scope": "openid offline_access",
        "redirect_uris": ["http://localhost:3000/"],
        "post_logout_redirect_uris": ["http://localhost:3000/"]
    }
    ```
- ## 3. Launch the identity provider

    To run the identity provider, execute the following at the root of the project.

    ```
    docker-compose up
    ```

- ## 4. Launch the app

    To run the identity provider, execute the following in this folder.

    ```
    yarn start
    ```

    Open [http://localhost:3000](http://localhost:3000) to view it in the browser.
