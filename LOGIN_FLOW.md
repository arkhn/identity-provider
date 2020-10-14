# Login and consent flow

The login flow implemented by Hydra is well documented here: https://www.ory.sh/hydra/docs/concepts/login.

## Client registration

First of all your client should be known by Hydra to perform any authentication flow.
To learn about how to register a client, please read

To perform a code authorization flow (the one we'll talk about here) your client will need to be registered with a grant type `authorization_code` and with response types `code` and `token`.

## How to use Arkhn's oauth provider in your web app

All of the queries presented below could be implemented by hand but it's preferable to use an already well-tested library that can do it for you. For instance, in our node projects, we use `client-oauth2` at Arkhn.

Also, please note that in the following code snippets, some helper functions are used to manipulate the tokens (as `getAccessToken`, `refreshToken`, etc.). An example of Arkhn's implementation for these can be found here: [tokenManager.ts](https://github.com/arkhn/warehouse-api/blob/master/front/src/services/tokenManager.ts).

The steps to perform the authentication flow on the application side are:

- ## Instantiating a ClientOAuth2

  As we said, we used a library called `client-oauth2` instead of implementing manually all the interaction between the different parties. Here is how to instantiate a client

  ```js
  export const authClient = new ClientOAuth2({
    clientId: CLIENT_ID,
    clientSecret: CLIENT_SECRET,
    authorizationUri: AUTH_URL,
    accessTokenUri: TOKEN_URL,
    redirectUri: LOGIN_REDIRECT_URL,
    scopes: ["openid", "offline_access"],
  });
  ```

  The important things to observe are the URL used to communicate with Hydra:

  - `authorizationUri`: the route to call to initiate a authorization code flow. It will be `<hydra-url>/oauth2/auth`.
  - `accessTokenUri`: the route to call to exchange a code for a token or to refresh a token. It will be `<hydra-url>/oauth2/token`.

  The other parameters are:

  - `clientId` and `clientSecret`: they identify and authenticate the application during the communications with Hydra.
  - `scopes`: the scopes that will be asked for to the user. They depend on what the appliation should be allowed to do on behalf of the user.
  - `redirectUri`: the uri to which the user should be redirected after it has logged in. It is basically the URL of your application's home page.

- ## Start the authentication flow:

  This basically consists in a GET query to Hydra's `/oauth2/auth` route.

  With `client-oauth2`:

  ```js
  const uri = authClient.code.getUri({
    state: state,
  });
  window.location.assign(uri);
  ```

- ## Fill in login

  The user will be redirected to the login and consent (optional) pages.

- ## Use the code given by Hydra

  At the end of the flow, the user is redirected to the application with a code. This code exchanged for a token.
  Once again, this could be done manually but it's easier and safer to use a library.

  You can use the access token in you queries by storing it in your localstorage or using cookies.

  With `client-oauth2`:

  ```js
  // oauthToken is an object that contains an access token,
  // a refresh token (optional), and additional data if relevant.
  const oauthToken = await authClient.code.getToken(window.location.href);
  // if you want to store the access token in your local storage
  localStorage.setItem(ACCESS_TOKEN_STORAGE_KEY, oauthToken.accessToken);
  ```

- ## Set the authorization header of your queries that need the access token

  If you have stored the token somewhere, you'll need to add it to the Authorization header in your queries.
  If you're using axios, it's easy to do with an interceptor:

  ```js
  // Interceptor to add authorization header for each requests to the API
  axios.interceptors.request.use((config) => {
    if (config.url?.startsWith(FHIR_API_URL)) {
      const token = getAccessToken();
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });
  ```

- ## Refresh your token

  Hydra's tokens expire after a while. If you don't want users to have to log in again after the token expiration, you can use a refresh token to automatically get a new valid access token.

  Once again, if you use axios, you can do that with an interceptor:

  ```js
  // Add an interceptor to refresh access token when needed
  axios.interceptors.response.use(
    (response) => {
      return response;
    },
    async (error) => {
      const originalRequest = error.config;

      if (
        error.response.status === 401 &&
        originalRequest.url.startsWith(TOKEN_URL)
      ) {
        redirectToLogin();
        return Promise.reject(error);
      }

      if (error.response.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;

        const success = await refreshToken();
        if (!success) {
          redirectToLogin();
          return Promise.reject(error);
        }
        return axios(originalRequest);
      }
      return Promise.reject(error);
    }
  );
  ```

  Please note that to be able to refresh a token, the scopes you ask for during the code authorization flow should include `offline_access` and your client should be registered with a grant_type `refresh_token`.

- ## Revoking tokens, logging out

  Token revokation is not implemented by `client-oauth2` so we had to do it manually. Here is what has to be done to revoke a token:

  ```js
  export const revokeToken = async () => {
    const accessToken = localStorage.getItem(ACCESS_TOKEN_STORAGE_KEY);
    if (!accessToken)
      throw new Error(
        "Access token not present in local storage, cannot revoke it."
      );

    const bodyFormData = new FormData();
    bodyFormData.set("token", accessToken);
    const conf = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
        Accept: "application/json, application/x-www-form-urlencoded",
        Authorization: "Basic " + btoa(CLIENT_ID + ":" + CLIENT_SECRET),
      },
    };
    try {
      const revokeResponse = await axios.post(REVOKE_URL, bodyFormData, conf);
      if (revokeResponse.status !== 200) console.error(revokeResponse.data);
    } catch (err) {
      console.error(err.response);
    }
  };
  ```

  Hydra also has a logout flow. We won't present it here.
